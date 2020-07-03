package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	vrsSdk "github.com/nuagenetworks/libvrsdk/api"
	"github.com/nuagenetworks/libvrsdk/api/port"
	"github.com/nuagenetworks/nuage-cni/client"
	"github.com/nuagenetworks/nuage-cni/config"
	"github.com/nuagenetworks/nuage-cni/k8s"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	kclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var interruptChannel chan bool
var signalChannel chan os.Signal
var staleEntityMap map[string]int64
var stalePortMap map[string]int64
var staleEntryTimeout int64
var isAtomic bool
var hostname string
var orchestratorType string

// filter on host name
const (
	PodHostField = "spec.nodeName"
)

// cleanupStaleEntities will clear stale
// entity entries from Nuage tables
func cleanupStaleEntities(vrsConnection vrsSdk.VRSConnection, orchestrator string) error {

	log.Debugf("Cleaning up stale ports and entities in VRS as a part of the audit daemon")
	var err error
	var portList []string
	var k8sActivePortList []string
	var k8sActivePodNames []string
	var vrsEntityNames []string

	// First obtain VRS entity and port list followed by orchestrator
	// entity/port list to avoid race condition
	vrsEntitiesList, err := vrsConnection.GetAllEntities()
	if err != nil {
		log.Errorf("Failed to get entity list from VRS: %v", err)
		return err
	}
	for _, entityID := range vrsEntitiesList {
		entityName, err := vrsConnection.GetEntityName(entityID)
		if err != nil {
			log.Debugf("Error obtaining entity name from OVSDB: %v", err)
		}
		vrsEntityNames = append(vrsEntityNames, entityName)
	}

	vrsPortsList, err := vrsConnection.GetAllPorts()
	if err != nil {
		log.Errorf("Failed getting port names from VRS: %v", err)
		return err
	}

	k8sEntityMap, err := getActiveK8SPods(orchestrator)
	if err != nil {
		log.Errorf("Error occured while obtaining currently active Pods list: %v", err)
		return err
	}
	log.Debugf("Currently active k8s pods mapping : %v", k8sEntityMap)
	for _, name := range k8sEntityMap {
		k8sActivePodNames = append(k8sActivePodNames, name)
		portList, err = vrsConnection.GetEntityPortsByName(name)
		if err != nil {
			log.Warnf("Error while obtaining VRS ports with runtime entity name %s: %v", name, err)
		}
		k8sActivePortList = append(k8sActivePortList, portList...)
	}

	err = cleanupVMTable(vrsConnection, vrsEntityNames, k8sActivePodNames)
	if err != nil {
		log.Warnf("Cleaning up VM table failed with error %v", err)
	}

	err = cleanupPortTable(vrsConnection, vrsPortsList, k8sActivePortList)
	if err != nil {
		log.Warnf("Cleaning up port table failed with error %v", err)
	}

	return err
}

// getStaleEntityEntriesForDeletion will determine what entity entries
// need to be actually cleaned up from VRS
func getStaleEntityEntriesForDeletion(ids []string) []string {

	var deleteEntityList []string
	currentTime := (time.Now().UnixNano()) / 1000000
	for _, staleID := range ids {
		if _, ok := staleEntityMap[staleID]; !ok {
			staleEntityMap[staleID] = currentTime
		}

		timeDiff := (currentTime - staleEntityMap[staleID]) / 1000
		if timeDiff >= staleEntryTimeout {
			deleteEntityList = append(deleteEntityList, staleID)
		}
	}

	// Delete resolved entities earlier marked as stale
	// from stale entity map
	keyFound := false
	for key := range staleEntityMap {
		for _, staleID := range ids {
			if key == staleID {
				log.Debugf("Entry %s is still not resolved or is a stale entry", key)
				keyFound = true
				break
			}
		}
		if !keyFound {
			delete(staleEntityMap, key)
		} else {
			keyFound = false
		}
	}

	return deleteEntityList
}

// getStalePortEntriesForDeletion will determine what port entries
// need to be actually cleaned up from VRS
func getStalePortEntriesForDeletion(ids []string) []string {

	var deletePortList []string

	currentTime := (time.Now().UnixNano()) / 1000000
	for _, staleID := range ids {
		if _, ok := stalePortMap[staleID]; !ok {
			stalePortMap[staleID] = currentTime
		}

		timeDiff := (currentTime - stalePortMap[staleID]) / 1000
		if timeDiff >= staleEntryTimeout {
			deletePortList = append(deletePortList, staleID)
		}
	}

	// Delete resolved alubr0 ports earlier marked as stale
	// from stale port map
	keyFound := false
	for key := range stalePortMap {
		for _, staleID := range ids {
			if key == staleID {
				log.Debugf("Entry %s is still not resolved or is a stale entry", key)
				keyFound = true
				break
			}
		}
		if !keyFound {
			delete(staleEntityMap, key)
		} else {
			keyFound = false
		}
	}

	return deletePortList
}

func auditEntity(vrsConnection vrsSdk.VRSConnection, id string) bool {

	vrsPortsList, err := vrsConnection.GetEntityPorts(id)
	if err != nil {
		log.Errorf("Failed getting port names from VRS: %v", err)
		return true
	}

	for _, port := range vrsPortsList {
		if strings.HasPrefix(port, "nu") {
			return true
		}
	}
	return false
}

// cleanupVMTable removes stale entity entries from Nuage VM table
func cleanupVMTable(vrsConnection vrsSdk.VRSConnection, vrsEntityNameList []string, k8sActivePodNames []string) error {

	var err error
	var deleteStaleEntitiesList []string
	log.Debugf("Cleaning up stale entity entries from Nuage VM table")
	staleNames := computeStaleEntitiesDiff(vrsEntityNameList, k8sActivePodNames)
	deleteStaleEntitiesList = getStaleEntityEntriesForDeletion(staleNames)
	for _, staleName := range deleteStaleEntitiesList {
		doAudit := auditEntity(vrsConnection, staleName)
		if doAudit {
			log.Infof("Removing stale entity entry %s", staleName)
			ports, err := vrsConnection.GetEntityPortsByName(staleName)
			if err != nil {
				log.Debugf("Failed getting port names from VRS: %v", err)
			}
			err = vrsConnection.DestroyEntityByVMName(staleName)
			if err != nil {
				log.Warnf("Unable to delete entry from nuage VM table: %v", err)
			} else {
				sendStaleEntryDeleteNotification(vrsConnection, staleName, ports)
			}
			delete(staleEntityMap, staleName)
		} else {
			log.Debugf("Skipping Nuage audit as this is not CNI created entity entry")
			return nil
		}
	}

	log.Infof("Stale entities cleaned up from VRS %v", deleteStaleEntitiesList)
	return err
}

// sendStaleEntryDeleteNotification notifies monitor about
// stale VRS entity and port entry deletion
func sendStaleEntryDeleteNotification(vrsConnection vrsSdk.VRSConnection, entityName string, ports []string) {

	var err error

	var portInfo map[port.StateKey]interface{}
	for _, port := range ports {
		portInfo, err = vrsConnection.GetPortState(port)
		if err != nil {
			log.Debugf("Unable to obtain port Nuage metadata from VRS")
		}
	}

	if _, ok := portInfo[port.StateKeyNuageZone].(string); ok {
		log.Debugf("Sending delete notification for entity %s for zone %s", entityName, portInfo[port.StateKeyNuageZone].(string))
		// Send pod deletion notification to Nuage monitor
		err = k8s.SendPodDeletionNotification(entityName, portInfo[port.StateKeyNuageZone].(string), orchestratorType)
		if err != nil {
			log.Errorf("Error occured while sending delete notification for pod %s", entityName)
		}
	}
}

// cleanupPortTable removes stale port entries from Nuage VM table
func cleanupPortTable(vrsConnection vrsSdk.VRSConnection, vrsPortsList []string, entityPortList []string) error {
	var err error
	var deleteStalePortsList []string
	log.Debugf("Cleaning up stale entity entries from Nuage VM table")
	stalePorts := computeStaleEntitiesDiff(vrsPortsList, entityPortList)
	deleteStalePortsList = getStalePortEntriesForDeletion(stalePorts)
	for _, stalePort := range deleteStalePortsList {
		if strings.HasPrefix(stalePort, "nu") {
			log.Infof("Removing stale port %s", stalePort)
			err = vrsConnection.DestroyPort(stalePort)
			if err != nil {
				log.Warnf("Unable to delete port from Nuage Port table: %v", err)
			}

			// Purging out the veth port from VRS alubr0
			err = vrsConnection.RemovePortFromAlubr0(stalePort)
			if err != nil {
				log.Warnf("Unable to delete veth port as part of cleanup from alubr0: %v", err)
			}

			err = client.DeleteVethPair(stalePort, "eth0")
			if err != nil {
				log.Warnf("Failed to clear veth ports from VRS: %v", err)
			}

			delete(stalePortMap, stalePort)
		} else {
			log.Debugf("Skipping Nuage audit as this is not CNI created port entry")
			return nil
		}
	}

	log.Infof("Stale ports cleaned up from VRS %v", deleteStalePortsList)
	return err
}

// computeStalePortsEntitiesDiff will help determine the
// stale ports and entities in Nuage tables
func computeStaleEntitiesDiff(vrsData, orchestratorData []string) []string {
	log.Debugf("Computing stale entities on agent node as a part of the monitoring daemon")
	lookup := make(map[string]int)
	var res []string

	for _, str := range orchestratorData {
		lookup[str]++
	}

	for _, str := range vrsData {
		if _, ok := lookup[str]; !ok {
			res = append(res, str)
		}
	}
	return res
}

// handleDaemonInterrupt will handle any external interrupts
// to audit daemon and handle stale connection/entities cleanup
// to have a graceful daemon exit
func handleDaemonInterrupt() {

	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		for {
			s := <-signalChannel
			switch s {
			case syscall.SIGHUP:
				log.Errorf("SIGHUP signal interrupted Nuage CNI daemon")
				interruptChannel <- true
			case syscall.SIGINT:
				log.Errorf("SIGINT signal interrupted Nuage CNI daemon")
				interruptChannel <- true
			case syscall.SIGTERM:
				log.Errorf("SIGTERM signal interrupted Nuage CNI daemon")
				interruptChannel <- true
			case syscall.SIGQUIT:
				log.Errorf("SIGHQUIT signal interrupted Nuage CNI daemon")
				interruptChannel <- true
			default:
				break
			}
		}
	}()
}

// MonitorAgent will be run as a background audit daemon
// on k8s agent nodes to clean up stale entities/ports
// on agent nodes
func MonitorAgent(config *config.Config, orchestrator string) error {

	var err error
	var vrsConnection vrsSdk.VRSConnection
	interruptChannel = make(chan bool)

	staleEntityMap = make(map[string]int64)
	stalePortMap = make(map[string]int64)
	staleEntryTimeout = config.StaleEntryTimeout
	orchestratorType = orchestrator

	hostname, err = os.Hostname()
	if err != nil {
		log.Errorf("finding hostname failed with error: %v", err)
		return err
	}
	for {
		vrsConnection, err = client.ConnectToVRSOVSDB(config)
		if err != nil {
			log.Errorf("Error connecting to VRS. Will re-try connection in 5 seconds")
		} else {
			break
		}
		time.Sleep(time.Duration(5) * time.Second)
	}

	log.Infof("Starting Nuage CNI monitoring daemon for %s node with hostname %s", orchestrator, hostname)

	// Cleaning up stale ports/entities when audit daemon starts
	err = cleanupStaleEntities(vrsConnection, orchestrator)
	if err != nil {
		log.Errorf("Error cleaning up stale entities and ports on VRS")
	}

	// Determine whether the base host is RHEL server or RHEL atomic
	isAtomic = k8s.VerifyHostType()

	if !isAtomic && orchestrator == "ose" {
		cmdstr := fmt.Sprintf("rm -irf /var/usr/")
		cmd := exec.Command("bash", "-c", cmdstr)
		_, _ = cmd.Output()
	}

	vrsStaleEntriesCleanupTicker := time.NewTicker(time.Duration(config.MonitorInterval) * time.Second)
	vrsConnectionCheckTicker := time.NewTicker(time.Duration(config.VRSConnectionCheckTimer) * time.Second)

	handleDaemonInterrupt()

	for {
		select {
		case <-vrsStaleEntriesCleanupTicker.C:
			err := cleanupStaleEntities(vrsConnection, orchestrator)
			if err != nil {
				log.Errorf("Error cleaning up stale entities and ports on VRS")
			}
		case <-vrsConnectionCheckTicker.C:
			_, err := vrsConnection.GetAllEntities()
			if err != nil {
				log.Errorf("VRS connection is down; will retry connection")
				vrsConnection, err = client.ConnectToVRSOVSDB(config)
				if err != nil {
					log.Errorf("Error connecting to VRS. Retry connection in %d seconds", config.VRSConnectionCheckTimer)
				} else {
					log.Infof("VRS connection is restored")
				}
			}
		case <-interruptChannel:
			log.Errorf("Daemon was interrupted by an external interrupt; will cleanup before exiting")
			vrsConnection.Disconnect()
			return fmt.Errorf("Daemon was interrupted by an external interrupt")
		}
	}
}

// getActiveK8SPods will help obtain UUID list
// for currently active pods on k8s cluster
func getActiveK8SPods(orchestrator string) (map[string]string, error) {

	log.Infof("Obtaining currently active K8S pods on agent node")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("creating the in-cluster config failed %v", err)
		return map[string]string{}, err
	}
	// creates the clientset
	kubeClient, err := kclient.NewForConfig(config)
	if err != nil {
		log.Errorf("Error trying to create kubeclient %v", err)
		return map[string]string{}, err
	}
	selector := fields.OneTermEqualSelector(PodHostField, hostname).String()
	listOpts := metav1.ListOptions{FieldSelector: selector}
	pods, err := kubeClient.CoreV1().Pods(metav1.NamespaceAll).List(listOpts)
	if err != nil {
		log.Errorf("Error occured while fetching pods from k8s api server")
		return map[string]string{}, err
	}

	entityMap := make(map[string]string)
	for _, pod := range pods.Items {
		entityMap[string(pod.UID)] = pod.Name
	}

	return entityMap, err
}
