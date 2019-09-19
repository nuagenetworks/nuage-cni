package daemon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	vrsSdk "github.com/nuagenetworks/libvrsdk/api"
	"github.com/nuagenetworks/libvrsdk/api/port"
	"github.com/nuagenetworks/nuage-cni/client"
	"github.com/nuagenetworks/nuage-cni/config"
	"github.com/nuagenetworks/nuage-cni/k8s"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var interruptChannel chan bool
var signalChannel chan os.Signal
var staleEntityMap map[string]int64
var stalePortMap map[string]int64
var staleEntryTimeout int64
var isAtomic bool
var orchestratorType string

type containerInfo struct {
	ID string `json:"container_id"`
}

// getActiveMesosContainers will return a list of currently
// active Mesos containers to help in audit cleanup
func getActiveMesosContainers() ([]string, error) {

	log.Infof("Obtaining currently active Mesos containers on agent node")
	var id []containerInfo
	var containerList []string

	name, err := os.Hostname()
	if err != nil {
		log.Errorf("Error reading hostname of the agent node")
		return containerList, err
	}
	url := "http://" + name + ":5051/containers"

	res, err := http.Get(url)
	if err != nil {
		log.Errorf("Error reading http endpoint response on mesos agent node")
		return containerList, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error parsing http endpoint response body on mesos agent node")
		return containerList, err
	}

	data := []byte(string(body))

	err = json.Unmarshal(data, &id)
	if err != nil {
		log.Errorf("Error unmarshalling http endpoint JSON response on mesos agent node")
		return containerList, err
	}

	if len(id) >= 1 {
		for index := range id {
			containerList = append(containerList, id[index].ID)
		}
	}

	return containerList, err
}

// cleanupStalePortsEntities will clear stale port or
// entity entries from Nuage tables
func cleanupStalePortsEntities(vrsConnection vrsSdk.VRSConnection, orchestrator string) error {

	log.Debugf("Cleaning up stale ports and entities in VRS as a part of the audit daemon")
	var err error
	var portList []string
	var entityPortList []string
	var entityUUIDList []string
	var formattedEntityUUIDList []string

	// First obtain VRS entity and port list followed by orchestrator
	// entity/port list to avoid race condition
	vrsEntitiesList, err := vrsConnection.GetAllEntities()
	if err != nil {
		log.Errorf("Failed to get entity list from VRS: %v", err)
		return err
	}

	vrsPortsList, err := vrsConnection.GetAllPorts()
	if err != nil {
		log.Errorf("Failed getting port names from VRS: %v", err)
		return err
	}

	switch orchestrator {
	case "mesos":
		entityUUIDList, err = getActiveMesosContainers()
		if err != nil {
			log.Errorf("Error occured while obtaining currently active container list: %v", err)
			return err
		}
		log.Debugf("Currently active Mesos containers list : %v", entityUUIDList)
	case "k8s":
		entityUUIDList, err = getActiveK8SPods(orchestrator)
		if err != nil {
			log.Errorf("Error occured while obtaining currently active Pods list: %v", err)
			return err
		}
		log.Debugf("Currently active k8s pods list : %v", entityUUIDList)
	case "ose":
		entityUUIDList, err = getActiveK8SPods(orchestrator)
		if err != nil {
			log.Errorf("Error occured while obtaining currently active Pods list: %v", err)
			return err
		}
		log.Debugf("Currently active openshift pods list : %v", entityUUIDList)
	default:
	}

	var formattedUUID string
	for _, id := range entityUUIDList {
		if orchestrator == "mesos" {
			newID := strings.Replace(id, "-", "", -1)
			formattedUUID = newID + newID
		} else {
			formattedUUID = id
		}
		formattedEntityUUIDList = append(formattedEntityUUIDList, formattedUUID)
		portList, err = vrsConnection.GetEntityPorts(formattedUUID)
		if err != nil {
			log.Warnf("Error occured while obtaining VRS ports for entity %s", formattedUUID)
		}
		entityPortList = append(entityPortList, portList...)
	}

	err = cleanupVMTable(vrsConnection, formattedEntityUUIDList, vrsEntitiesList)
	if err != nil {
		log.Warnf("Cleaning up VM table failed with error %v", err)
	}

	err = cleanupPortTable(vrsConnection, entityPortList, vrsPortsList)
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
	for key, _ := range staleEntityMap {
		for _, staleID := range ids {
			if key == staleID {
				log.Debugf("Entry %s is still not resolved or is a stale entry", key)
				keyFound = true
				break
			}
		}
		if keyFound == false {
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
	for key, _ := range stalePortMap {
		for _, staleID := range ids {
			if key == staleID {
				log.Debugf("Entry %s is still not resolved or is a stale entry", key)
				keyFound = true
				break
			}
		}
		if keyFound == false {
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
func cleanupVMTable(vrsConnection vrsSdk.VRSConnection, entityUUIDList []string, vrsEntitiesList []string) error {

	var err error
	var deleteStaleEntitiesList []string
	log.Debugf("Cleaning up stale entity entries from Nuage VM table")
	staleIDs := computeStalePortsEntitiesDiff(vrsEntitiesList, entityUUIDList)
	deleteStaleEntitiesList = getStaleEntityEntriesForDeletion(staleIDs)
	for _, staleID := range deleteStaleEntitiesList {
		doAudit := auditEntity(vrsConnection, staleID)
		if doAudit == true {
			entityName, err := vrsConnection.GetEntityName(staleID)
			if err != nil {
				log.Debugf("Error obtaining entity name from OVSDB: %v", err)
			}

			ports, err := vrsConnection.GetEntityPorts(staleID)
			if err != nil {
				log.Debugf("Failed getting port names from VRS: %v", err)
			}

			log.Infof("Removing stale entity entry %s", staleID)
			err = vrsConnection.DestroyEntity(staleID)
			if err != nil {
				log.Warnf("Unable to delete entry from nuage VM table: %v", err)
			} else {
				sendStaleEntryDeleteNotification(vrsConnection, entityName, ports)
			}
			delete(staleEntityMap, staleID)
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
func cleanupPortTable(vrsConnection vrsSdk.VRSConnection, entityPortList []string, vrsPortsList []string) error {

	var err error
	var deleteStalePortsList []string
	log.Debugf("Cleaning up stale entity entries from Nuage VM table")
	stalePorts := computeStalePortsEntitiesDiff(vrsPortsList, entityPortList)
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
func computeStalePortsEntitiesDiff(vrsData, orchestratorData []string) []string {
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
// on Mesos/k8s agent nodes to clean up stale entities/ports
// on agent nodes
func MonitorAgent(config *config.Config, orchestrator string) error {

	var err error
	var vrsConnection vrsSdk.VRSConnection
	interruptChannel = make(chan bool)

	staleEntityMap = make(map[string]int64)
	stalePortMap = make(map[string]int64)
	staleEntryTimeout = config.StaleEntryTimeout
	orchestratorType = orchestrator

	for {
		vrsConnection, err = client.ConnectToVRSOVSDB(config)
		if err != nil {
			log.Errorf("Error connecting to VRS. Will re-try connection in 5 seconds")
		} else {
			break
		}
		time.Sleep(time.Duration(5) * time.Second)
	}

	log.Infof("Starting Nuage CNI monitoring daemon for %s agent nodes", orchestrator)

	// Cleaning up stale ports/entities when audit daemon starts
	err = cleanupStalePortsEntities(vrsConnection, orchestrator)
	if err != nil {
		log.Errorf("Error cleaning up stale entities and ports on VRS")
	}

	// Determine whether the base host is RHEL server or RHEL atomic
	isAtomic = k8s.VerifyHostType()

	if isAtomic == false && orchestrator == "ose" {
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
			err := cleanupStalePortsEntities(vrsConnection, orchestrator)
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
func getActiveK8SPods(orchestrator string) ([]string, error) {

	log.Infof("Obtaining currently active K8S pods on agent node")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("creating the in-cluster config failed %v", err)
		return []string{}, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("creating clientset for k8s failed %v", err)
		return []string{}, err
	}

	var listOpts = metav1.ListOptions{}
	pods, err := clientset.CoreV1().Pods("").List(listOpts)
	if err != nil {
		log.Errorf("Error occured while fetching pods from k8s api server")
		return []string{}, err
	}

	var ids []string
	for _, pod := range pods.Items {
		ids = append(ids, string(pod.UID))
	}

	return ids, err
}
