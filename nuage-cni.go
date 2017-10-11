//
// Copyright (c) 2016 Nuage Networks, Inc. All rights reserved.
//

// This will form the Nuage CNI plugin for networking
// containers spawned using Mesos & k8s containerizers

package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/version"
	vrsSdk "github.com/nuagenetworks/libvrsdk/api"
	"github.com/nuagenetworks/libvrsdk/api/entity"
	"github.com/nuagenetworks/libvrsdk/api/port"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"nuage-cni/client"
	"nuage-cni/config"
	"nuage-cni/daemon"
	"nuage-cni/k8s"
	"os"
	"runtime"
	"strings"
	"time"
)

var hostname string
var logMessageCounter int

type logTextFormatter log.TextFormatter

var supportedLogLevels = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
}

// nuageCNIConfig will be a pointer variable to
// Config struct that will hold all Nuage CNI plugin
// parameters
var nuageCNIConfig = &config.Config{}

var operMode string
var orchestrator string

// nuageMetadataObj will be a structure pointer
// to hold Nuage metadata
var nuageMetadataObj = client.NuageMetadata{}

// Const definitions for plugin log location and input parameter file
const (
	paramFile     = "/etc/default/nuage-cni.yaml"
	logFolder     = "/var/log/cni/"
	cniLogFile    = "/var/log/cni/nuage-cni.log"
	daemonLogFile = "/var/log/cni/nuage-daemon.log"
	bridgeName    = "alubr0"
	kubernetes    = "k8s"
	mesos         = "mesos"
	openshift     = "ose"
)

func init() {
	// This ensures that main runs only on main thread (thread group leader).
	// since namespace ops (unshare, setns) are done for a single thread, we
	// must ensure that the goroutine does not jump from OS thread to thread
	runtime.LockOSThread()

	hostname, _ = os.Hostname()

	// Reading Nuage CNI plugin parameter file
	data, err := ioutil.ReadFile(paramFile)
	if err != nil {
		log.Errorf("Error in reading from Nuage CNI plugin parameter file: %s\n", err)
	}

	if err = yaml.Unmarshal(data, nuageCNIConfig); err != nil {
		log.Errorf("Error in unmarshalling data from Nuage CNI parameter file: %s\n", err)
	}

	// Set default values if some values were not set
	// in Nuage CNI yaml file
	client.SetDefaultsForNuageCNIConfig(nuageCNIConfig)

	if _, err = os.Stat(logFolder); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(logFolder, 777)
			if err != nil {
				fmt.Printf("Error creating log folder: %v", err)
			}
		}
	}

	// Use a new flag set so as not to conflict with existing
	// libraries which use "flag"
	flagSet := flag.NewFlagSet("Nuage", flag.ExitOnError)

	// Determining the mode of operation
	mode := flagSet.Bool("daemon", false, "a bool")

	err = flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var logfile string
	if *mode {
		operMode = "daemon"
		logfile = daemonLogFile
	} else {
		operMode = "cni"
		logfile = cniLogFile
	}

	customFormatter := new(logTextFormatter)
	log.SetFormatter(customFormatter)
	log.SetOutput(&lumberjack.Logger{
		Filename: logfile,
		MaxSize:  nuageCNIConfig.LogFileSize,
		MaxAge:   30,
	})
	log.SetLevel(supportedLogLevels[strings.ToLower(nuageCNIConfig.LogLevel)])

	// Determine which orchestrator is making the CNI call
	var arg string
	arg = os.Args[0]
	if strings.Contains(arg, mesos) {
		orchestrator = mesos
	} else if strings.Contains(arg, kubernetes) {
		orchestrator = kubernetes
	} else {
		orchestrator = openshift
	}

	switch orchestrator {
	case mesos:
		log.Debugf("CNI call for mesos orchestrator")
	case kubernetes:
		log.Debugf("CNI call for k8s orchestrator")
	case openshift:
		log.Debugf("CNI call for ose orchestrator")
	default:
		panic("Invalid orchestrator for the CNI call")
	}
}

func (f *logTextFormatter) Format(entry *log.Entry) ([]byte, error) {
	logMessageCounter++
	return []byte(fmt.Sprintf("|%v|%s|%04d|%s\n", entry.Time, strings.ToUpper(log.Level.String(entry.Level)), logMessageCounter, entry.Message)), nil
}

func networkConnect(args *skel.CmdArgs) error {

	log.Infof("Nuage CNI plugin invoked to add an entity to Nuage defined VSD network")
	var err error
	var vrsConnection vrsSdk.VRSConnection
	var result *types.Result
	entityInfo := make(map[string]string)

	for {
		vrsConnection, err = client.ConnectToVRSOVSDB(nuageCNIConfig)
		if err != nil {
			log.Errorf("Error connecting to VRS. Will re-try connection")
		} else {
			break
		}
		time.Sleep(time.Duration(3) * time.Second)
	}
	log.Debugf("Successfully established a connection to Nuage VRS")

	if orchestrator == kubernetes || orchestrator == openshift {
		log.Debugf("Orchestrator ID is %s", orchestrator)
		// Parsing CNI args obtained for K8S/Openshift
		k8sArgs := client.K8sArgs{}
		err = types.LoadArgs(args.Args, &k8sArgs)
		if err != nil {
			log.Errorf("Error in loading k8s CNI arguments")
			return fmt.Errorf("Error in loading k8s CNI arguments: %s", err)
		}

		log.Debugf("Infra Container ID for pod %s is %s", string(k8sArgs.K8S_POD_NAME), string(k8sArgs.K8S_POD_INFRA_CONTAINER_ID))

		entityInfo["name"] = string(k8sArgs.K8S_POD_NAME)
		entityInfo["entityport"] = args.IfName
		entityInfo["brport"] = client.GetNuagePortName(entityInfo["entityport"], args.ContainerID)
		err := k8s.GetPodNuageMetadata(&nuageMetadataObj, string(k8sArgs.K8S_POD_NAME), string(k8sArgs.K8S_POD_NAMESPACE), orchestrator)
		if err != nil {
			log.Errorf("Error obtaining Nuage metadata")
			return fmt.Errorf("Error obtaining Nuage metadata: %s", err)
		}
		entityInfo["uuid"] = string(k8sArgs.K8S_POD_INFRA_CONTAINER_ID)
		log.Infof("Nuage metadata obtained for pod %s is Enterprise: %s, Domain: %s, Zone: %s, Network: %s and User:%s", string(k8sArgs.K8S_POD_NAME), nuageMetadataObj.Enterprise, nuageMetadataObj.Domain, nuageMetadataObj.Zone, nuageMetadataObj.Network, nuageMetadataObj.User)
	} else {
		log.Debugf("Orchestrator ID is %s", orchestrator)
		entityInfo["name"] = args.ContainerID
		newContainerUUID := strings.Replace(args.ContainerID, "-", "", -1)
		formattedContainerUUID := newContainerUUID + newContainerUUID
		entityInfo["uuid"] = formattedContainerUUID
		entityInfo["entityport"] = args.IfName
		entityInfo["brport"] = client.GetNuagePortName(entityInfo["entityport"], entityInfo["uuid"])
		err = client.GetContainerNuageMetadata(&nuageMetadataObj, args)
		if err != nil {
			log.Errorf("Error obtaining Nuage metadata")
			return fmt.Errorf("Error obtaining Nuage metadata: %s", err)
		}
	}

	// Verifying all required Nuage metadata present before proceeding
	if (nuageMetadataObj.Enterprise == "") || (nuageMetadataObj.Domain == "") || (nuageMetadataObj.Zone == "") || (nuageMetadataObj.Network == "") || (nuageMetadataObj.User == "") {
		log.Errorf("Required Nuage metadata not available for port resolution")
		return fmt.Errorf("Required Nuage metadata not available for port resolution")
	}

	log.Infof("Attaching entity %s to Nuage defined network", entityInfo["name"])

	// Here we setup veth paired interface to connect the Container
	// to Nuage defined network
	netns := args.Netns
	contVethMAC, err := client.SetupVEth(netns, entityInfo, nuageCNIConfig.MTU)
	if err != nil {
		log.Errorf("Error creating veth paired interface for entity %s", entityInfo["name"])
		// Cleaning up veth ports from VRS before returning if we fail during
		// veth create task
		_ = client.DeleteVethPair(entityInfo["brport"], entityInfo["entityport"])
		return fmt.Errorf("Failed to create veth paired interface for the entity")
	}
	log.Debugf("Successfully created a veth paired port for entity %s", entityInfo["name"])

	var info vrsSdk.EntityInfo
	info.Name = entityInfo["name"]
	info.UUID = entityInfo["uuid"]
	err = vrsConnection.AddPortToAlubr0(entityInfo["brport"], info)
	if err != nil {
		log.Errorf("Error adding bridge veth end %s of entity %s to alubr0", entityInfo["brport"], entityInfo["name"])
		// Cleaning up veth ports from VRS
		_ = client.DeleteVethPair(entityInfo["brport"], entityInfo["entityport"])
		return fmt.Errorf("Failed to add bridge veth port to alubr0")
	}
	log.Debugf("Attached veth interface %s to bridge %s for entity %s", entityInfo["brport"], bridgeName, entityInfo["name"])

	// Create Port Attributes
	portAttributes := port.Attributes{
		Platform: entity.Container,
		MAC:      contVethMAC,
		Bridge:   bridgeName,
	}

	// Create Port Metadata for entity
	portMetadata := make(map[port.MetadataKey]string)
	portMetadata[port.MetadataKeyDomain] = nuageMetadataObj.Domain
	portMetadata[port.MetadataKeyNetwork] = nuageMetadataObj.Network
	portMetadata[port.MetadataKeyZone] = nuageMetadataObj.Zone
	portMetadata[port.MetadataKeyNetworkType] = "ipv4"

	// Handling static IP scenario
	if nuageMetadataObj.StaticIP != "" {
		portMetadata[port.MetadataKeyStaticIP] = nuageMetadataObj.StaticIP
	}

	// Handling policy group assignment scenario
	if nuageMetadataObj.PolicyGroup != "" {
		portMetadata[port.MetadataNuagePolicyGroup] = nuageMetadataObj.PolicyGroup
	}

	// Handling redirection target scenario
	if nuageMetadataObj.RedirectionTarget != "" {
		portMetadata[port.MetadataKeyNuageRedirectionTarget] = nuageMetadataObj.RedirectionTarget
	}

	// Create an entry for entity in Nuage Port Table
	err = vrsConnection.CreatePort(entityInfo["brport"], portAttributes, portMetadata)
	if err != nil {
		log.Errorf("Error creating entity port for entity %s in Nuage Port table", entityInfo["name"])
		_ = client.DeleteVethPair(entityInfo["brport"], entityInfo["entityport"])
		_ = vrsConnection.RemovePortFromAlubr0(entityInfo["brport"])
		return fmt.Errorf("Unable to create entity port %v", err)
	}
	log.Debugf("Successfully created a port for entity %s in Nuage Port table", entityInfo["name"])

	entityExists, _ := vrsConnection.CheckEntityExists(entityInfo["uuid"])
	if !entityExists {
		// Populate entity metadata
		entityMetadata := make(map[entity.MetadataKey]string)
		entityMetadata[entity.MetadataKeyUser] = nuageMetadataObj.User
		entityMetadata[entity.MetadataKeyEnterprise] = nuageMetadataObj.Enterprise

		// Define ports associated with the entity
		ports := []string{entityInfo["brport"]}

		// Add entity to VRS
		entityInfoVRS := vrsSdk.EntityInfo{
			UUID:     entityInfo["uuid"],
			Name:     entityInfo["name"],
			Domain:   entity.Docker,
			Type:     entity.Container,
			Ports:    ports,
			Metadata: entityMetadata,
		}

		// Sending proper events for container activation
		// as these are different for each entity type in VRS
		events := &entity.EntityEvents{}
		events.EntityEventCategory = entity.EventCategoryStarted
		events.EntityEventType = entity.EventStartedBooted
		events.EntityState = entity.Running
		events.EntityReason = entity.RunningBooted
		entityInfoVRS.Events = events

		err = vrsConnection.CreateEntity(entityInfoVRS)
		if err != nil {
			log.Errorf("Error creating an entry in Nuage entity table for entity %s", entityInfo["name"])
			return fmt.Errorf("Unable to add entity to VRS %v", err)
		}
		log.Debugf("Successfully created an entity in Nuage entity table for entity %s", entityInfo["name"])
	} else {
		portList, err := vrsConnection.GetEntityPorts(entityInfo["uuid"])
		if err != nil {
			log.Errorf("Error obtaining alubr0 port for entity %s", entityInfo["name"])
			return fmt.Errorf("Unable to obtain alubr0 port for entity %v", err)
		}
		if len(portList) > 0 {
			entityInfo["brport"] = portList[0]
			log.Infof("Using existing alubr0 port %s for entity %s", entityInfo["brport"], entityInfo["name"])
		} else {
			log.Errorf("Error configuring alubr0 port for an existing VRS entity %s", entityInfo["name"])
			return fmt.Errorf("Error configuring alubr0 port for an existing VRS entity %v", err)
		}
	}

	// Registering for VRS port updates
	portInfoUpdateChan := make(chan *vrsSdk.PortIPv4Info)
	err = vrsConnection.RegisterForPortUpdates(entityInfo["brport"], portInfoUpdateChan)
	if err != nil {
		log.Errorf("Failed to register for updates from VRS for entity port %s", entityInfo["brport"])
		return fmt.Errorf("Failed to register for updates from VRS %v", err)
	}
	ticker := time.NewTicker(time.Duration(nuageCNIConfig.PortResolveTimer) * time.Second)
	portInfo := &vrsSdk.PortIPv4Info{}
	select {
	case portInfo = <-portInfoUpdateChan:
		log.Debugf("Received an update from VRS for entity port %s", entityInfo["brport"])
	case <-ticker.C:
		log.Errorf("Failed to receive an update from VRS for entity port %s", entityInfo["brport"])
		return fmt.Errorf("Failed to receive an IP address from Nuage CNI plugin%v", err)
	}

	// Configuring entity end veth with IP
	entityInfo["ip"] = portInfo.IPAddr
	entityInfo["gw"] = portInfo.Gateway
	entityInfo["mask"] = portInfo.Mask

	result, err = client.AssignIPToContainerIntf(netns, entityInfo)
	if err != nil {
		log.Errorf("Error configuring entity %s with an IP address", entityInfo["name"])
		return fmt.Errorf("Error configuring entity interface with IP %v", err)
	}

	log.Infof("Successfully configured entity %s with an IP address %s", entityInfo["name"], entityInfo["ip"])

	// De-registering for VRS port updates
	err = vrsConnection.DeregisterForPortUpdates(entityInfo["brport"])
	if err != nil {
		log.Errorf("Error de-registering for port updates from VRS for entity port %s", entityInfo["brport"])
	}

	vrsConnection.Disconnect()

	return result.Print()
}

func networkDisconnect(args *skel.CmdArgs) error {

	log.Infof("Nuage CNI plugin invoked to detach an entity from a Nuage defined VSD network")
	var err error
	var vrsConnection vrsSdk.VRSConnection
	var portName string
	entityInfo := make(map[string]string)

	if orchestrator == kubernetes || orchestrator == openshift {
		log.Debugf("Orchestrator ID is %s", orchestrator)
		// Parsing CNI args obtained for K8S
		k8sArgs := client.K8sArgs{}
		err = types.LoadArgs(args.Args, &k8sArgs)
		if err != nil {
			log.Errorf("Error in loading k8s CNI arguments")
			return fmt.Errorf("Error in loading k8s CNI arguments: %s", err)
		}

		entityInfo["name"] = string(k8sArgs.K8S_POD_NAME)
		entityInfo["uuid"] = string(k8sArgs.K8S_POD_INFRA_CONTAINER_ID)
		entityInfo["entityport"] = args.IfName
		entityInfo["zone"] = string(k8sArgs.K8S_POD_NAMESPACE)
		// Determining the Nuage host port name to be deleted from OVSDB table
		portName = client.GetNuagePortName(args.IfName, args.ContainerID)
	} else {
		entityInfo["name"] = args.ContainerID
		newContainerUUID := strings.Replace(args.ContainerID, "-", "", -1)
		formattedContainerUUID := newContainerUUID + newContainerUUID
		entityInfo["uuid"] = formattedContainerUUID
		entityInfo["entityport"] = args.IfName
		// Determining the Nuage host port name to be deleted from OVSDB table
		portName = client.GetNuagePortName(args.IfName, entityInfo["uuid"])
	}

	log.Infof("Detaching entity %s from Nuage defined network", entityInfo["name"])

	for {
		vrsConnection, err = client.ConnectToVRSOVSDB(nuageCNIConfig)
		if err != nil {
			log.Errorf("Error connecting to VRS. Will re-try connection")
		} else {
			break
		}
		time.Sleep(time.Duration(3) * time.Second)
	}
	log.Debugf("Successfully established a connection to Nuage VRS")

	// Obtaining all ports associated with this entity
	portList, _ := vrsConnection.GetEntityPorts(entityInfo["uuid"])

	// Delete VRS OVSDB entries only if the ports for the entity
	// exist in VRS tables
	if len(portList) == 1 {
		// Performing cleanup of port/entity on VRS
		err = vrsConnection.DestroyPort(portName)
		if err != nil {
			log.Errorf("Failed to delete entity port from Nuage Port table for entity %s", entityInfo["name"])
		} else {
			log.Infof("Successfully deleted entity port and sending delete event to monitor")
			// Send pod deletion notification to Nuage monitor only if port deletion
			// in VRS succeeds
			err = k8s.SendPodDeletionNotification(entityInfo["name"], entityInfo["zone"], orchestrator)
			if err != nil {
				log.Errorf("Error occured while sending delete notification for pod %s", entityInfo["name"])
			}
		}

		// Purging out the veth port from VRS alubr0
		err = vrsConnection.RemovePortFromAlubr0(portName)
		if err != nil {
			log.Errorf("Failed to remove veth port %s for entity %s from alubr0", portName, entityInfo["name"])
		}

		// Cleaning up veth paired ports from VRS
		err = client.DeleteVethPair(portName, entityInfo["entityport"])
		if err != nil {
			log.Errorf("Failed to clear veth ports from VRS for entity %s", entityInfo["name"])
		}
	}

	// Check if entity exists in OVSDB before trying deletion
	entityExists, _ := vrsConnection.CheckEntityExists(entityInfo["uuid"])
	if entityExists {
		err = vrsConnection.DestroyEntity(entityInfo["uuid"])
		if err != nil {
			log.Errorf("Failed to remove entity from Nuage entity Table for entity %s", entityInfo["name"])
		}
	}

	vrsConnection.Disconnect()

	return nil
}

func main() {

	// This is added to handle https://github.com/kubernetes/kubernetes/pull/24983
	// which is a known k8s CNI issue
	if orchestrator == kubernetes || orchestrator == openshift {
		if err := client.AddIgnoreUnknownArgs(); err != nil {
			os.Exit(1)
		}
	}

	var err error
	if operMode == "daemon" {
		log.Infof("Starting Nuage CNI audit daemon on agent nodes")
		err = daemon.MonitorAgent(nuageCNIConfig, orchestrator)
		if err != nil {
			log.Errorf("Error encountered while running Nuage CNI daemon: %s\n", err)
		}
	} else {
		skel.PluginMain(networkConnect, networkDisconnect, version.PluginSupports("0.2.0", "0.3.0"))
	}
}
