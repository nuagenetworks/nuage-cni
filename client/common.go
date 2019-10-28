// This module includes utilities that can used by Nuage CNI plugin
// during runtime

package client

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/containernetworking/cni/pkg/ip"
	"github.com/containernetworking/cni/pkg/ns"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"
	vrsSdk "github.com/nuagenetworks/libvrsdk/api"
	"github.com/nuagenetworks/nuage-cni/config"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

const (
	addCli    = "add"
	deleteCli = "del"
)

// AddIgnoreUnknownArgs appends the 'IgnoreUnknown=1' option to CNI_ARGS before calling the CNI plugin.
// Otherwise, it will complain about the Kubernetes arguments.
// See https://github.com/kubernetes/kubernetes/pull/24983
func AddIgnoreUnknownArgs() error {
	cniArgs := "IgnoreUnknown=1"
	if os.Getenv("CNI_ARGS") != "" {
		cniArgs = fmt.Sprintf("%s;%s", cniArgs, os.Getenv("CNI_ARGS"))
	}
	return os.Setenv("CNI_ARGS", cniArgs)
}

// SetupVEth will set up veth pair for container
// to be connected to Nuage defined VSD network
func SetupVEth(netns string, containerInfo map[string]string, mtu int) (contVethMAC string, err error) {

	log.Debugf("Creating veth paired ports for container %s with container port %s and host port %s", containerInfo["name"], containerInfo["entityport"], containerInfo["brport"])

	// Creating veth pair to attach container ns to a host ns to
	// a Nuage network
	err = ns.WithNetNSPath(netns, func(hostNS ns.NetNS) error {
		localVethPair := &netlink.Veth{
			LinkAttrs: netlink.LinkAttrs{Name: containerInfo["brport"], MTU: mtu},
			PeerName:  containerInfo["entityport"],
		}

		if err = netlink.LinkAdd(localVethPair); err != nil {
			log.Errorf("Failed to create veth paired port for container %s", containerInfo["name"])
			return fmt.Errorf("Failed to create a veth paired port for the container")
		}

		contVeth, errStr := netlink.LinkByName(containerInfo["entityport"])
		if errStr != nil {
			log.Errorf("Failed to lookup container port %s for container %s", containerInfo["entityport"], containerInfo["name"])
			return fmt.Errorf("Failed to lookup container end veth port")
		}
		contVethMAC = contVeth.Attrs().HardwareAddr.String()
		log.Debugf("MAC address for container end of veth paired port for container %s is %s", containerInfo["name"], contVethMAC)

		// Bring the container side veth port up
		err = netlink.LinkSetUp(contVeth)
		if err != nil {
			log.Errorf("Error enabling container end of veth port for container %s", containerInfo["name"])
			return fmt.Errorf("Failed to container end veth interface")
		}

		brVeth, errStr := netlink.LinkByName(containerInfo["brport"])
		if errStr != nil {
			log.Errorf("Failed to lookup host port %s for container %s", containerInfo["brport"], containerInfo["name"])
			return fmt.Errorf("Failed to lookup alubr0 bridge end veth port")
		}

		// Now that the everything has been successfully set up in the container, move the "host" end of the
		// veth into the host namespace.
		if err = netlink.LinkSetNsFd(brVeth, int(hostNS.Fd())); err != nil {
			return fmt.Errorf("failed to move veth to host netns: %v", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	brVeth, err := netlink.LinkByName(containerInfo["brport"])
	if err != nil {
		log.Errorf("Failed to lookup host port %s for container %s", containerInfo["brport"], containerInfo["name"])
		return "", fmt.Errorf("failed to lookup %q: %v", containerInfo["brport"], err)
	}

	if err = netlink.LinkSetUp(brVeth); err != nil {
		log.Errorf("Error enabling host end of veth port for container %s", containerInfo["name"])
		return "", fmt.Errorf("failed to set %q up: %v", containerInfo["brport"], err)
	}

	return contVethMAC, err
}

// AssignIPToContainerIntf will configure the container end of the veth
// interface with IP address assigned by the Nuage CNI plugin
func AssignIPToContainerIntf(netns string, containerInfo map[string]string) (*current.Result, error) {

	var err error
	r := &current.Result{}

	log.Debugf("Configuring container %s interface %s with IP %s and default gateway %s assigned by Nuage CNI plugin", containerInfo["name"], containerInfo["entityport"], containerInfo["ip"], containerInfo["gw"])

	netmask := net.IPMask(net.ParseIP(containerInfo["mask"]).To4())
	prefixSize, _ := netmask.Size()

	ipV4Network := net.IPNet{IP: net.ParseIP(containerInfo["ip"]), Mask: net.CIDRMask(prefixSize, 32)}
	r.IPs = []*current.IPConfig{&current.IPConfig{Address: ipV4Network}}
	err = ns.WithNetNSPath(netns, func(hostNS ns.NetNS) error {

		contVeth, errStr := netlink.LinkByName(containerInfo["entityport"])
		if errStr != nil {
			err = fmt.Errorf("failed to lookup %q: %v", containerInfo["entityport"], err)
			return err
		}

		// Add a connected route to a dummy next hop so that a default route can be set
		gw := net.ParseIP(containerInfo["gw"])
		gwNet := &net.IPNet{IP: gw, Mask: net.CIDRMask(32, 32)}
		if err = netlink.RouteAdd(&netlink.Route{
			LinkIndex: contVeth.Attrs().Index,
			Scope:     netlink.SCOPE_LINK,
			Dst:       gwNet}); err != nil {
			return fmt.Errorf("failed to add route %v", err)
		}

		if err = ip.AddDefaultRoute(gw, contVeth); err != nil {
			log.Infof("Default route already exists within the container; skip re-configuring default route")
		} else {
			log.Debugf("Successfully added default route to container %s via gateway %s", containerInfo["name"], containerInfo["gw"])
		}
		for _, ip := range r.IPs {
			if err = netlink.AddrAdd(contVeth, &netlink.Addr{IPNet: &ip.Address}); err != nil {
				log.Errorf("Failed to assign IP %s to container %s", ip.Address, containerInfo["name"])
				return fmt.Errorf("failed to add IP addr to %q: %v", containerInfo["entityport"], err)
			}
			log.Debugf("Successfully assigned IP %s to container %s", ip.Address, containerInfo["name"])
		}

		return err
	})

	return r, err
}

// ConnectToVRSOVSDB will try connecting to VRS OVSDB via unix socket
// connection
func ConnectToVRSOVSDB(conf *config.Config) (vrsSdk.VRSConnection, error) {

	vrsConnection, err := vrsSdk.NewUnixSocketConnection(conf.VRSEndpoint)
	if err != nil {
		return vrsConnection, fmt.Errorf("Couldn't connect to VRS: %s", err)
	}

	return vrsConnection, nil
}

// DeleteVethPair will help user delete veth pairs on VRS
func DeleteVethPair(brPort string, entityPort string) error {

	log.Debugf("Deleting veth paired port %s as a part of Nuage CNI cleanup", brPort)
	localVethPair := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: brPort},
		PeerName:  entityPort,
	}

	err := netlink.LinkDel(localVethPair)
	if err != nil {
		log.Errorf("Deleting veth pair %+v failed with error: %s", localVethPair, err)
		return err
	}

	return nil
}

// GetNuagePortName creates a unique port name
// for container host port entry
func GetNuagePortName(uuid string) string {

	// Formatting UUID string by removing "-"
	formattedUUID := strings.Replace(uuid, "-", "", -1)
	nuagePortName := "nu" + generateVEthString(formattedUUID)

	return nuagePortName
}

// generateVEthString generates a unique SHA encoded
// string for veth Nuage host port
func generateVEthString(uuid string) string {
	h := sha1.New()
	_, err := h.Write([]byte(uuid))
	if err != nil {
		log.Errorf("Error generating unique hash string for entity")
	}
	return fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil))[:13])
}

// GetContainerNuageMetadata populates NuageMetadata struct
// with network information from labels passed from CNI NetworkProtobuf
func GetContainerNuageMetadata(nuageMetadata *NuageMetadata, args *skel.CmdArgs) error {

	var err error

	// Loading CNI network configuration
	conf := NetConf{}
	if err = json.Unmarshal(args.StdinData, &conf); err != nil {
		return fmt.Errorf("Failed to load netconf from CNI: %v", err)
	}

	// Parse endpoint labels passed in by Mesos, and store in a map.
	labels := map[string]string{}
	for _, label := range conf.Args.Mesos.NetworkInfo.Labels.Labels {
		labels[label.Key] = label.Value
	}

	if _, ok := labels["enterprise"]; ok {
		nuageMetadata.Enterprise = labels["enterprise"]
	}

	if _, ok := labels["domain"]; ok {
		nuageMetadata.Domain = labels["domain"]
	}

	if _, ok := labels["zone"]; ok {
		nuageMetadata.Zone = labels["zone"]
	}

	if _, ok := labels["network"]; ok {
		nuageMetadata.Network = labels["network"]
	}

	if _, ok := labels["user"]; ok {
		nuageMetadata.User = labels["user"]
	}

	if _, ok := labels["policy_group"]; ok {
		nuageMetadata.PolicyGroup = labels["policy_group"]
	}

	if _, ok := labels["static_ip"]; ok {
		nuageMetadata.StaticIP = labels["static_ip"]
	}

	if _, ok := labels["redirection_target"]; ok {
		nuageMetadata.RedirectionTarget = labels["redirection_target"]
	}

	return err
}

// SetDefaultsForNuageCNIConfig will set default values for
// Nuage CNI yaml parameters if they have not been set
func SetDefaultsForNuageCNIConfig(conf *config.Config) {

	if conf.VRSEndpoint == "" {
		log.Warnf("VRS endpoint not set. Using default value")
		conf.VRSEndpoint = "/var/run/openvswitch/db.sock"
	}

	if conf.VRSBridge == "" {
		log.Warnf("VRS bridge not set. Using default value")
		conf.VRSBridge = "alubr0"
	}

	if conf.MonitorInterval == 0 {
		log.Warnf("Monitor interval for audit daemon not set. Using default value")
		conf.MonitorInterval = 60
	}

	if conf.CNIVersion == "" {
		log.Warnf("CNI version not set. Using default value")
		conf.CNIVersion = "0.2.0"
	}

	if conf.PortResolveTimer == 0 {
		log.Warnf("OVSDB port resolution wait timer not set. Using default value")
		conf.PortResolveTimer = 60
	}

	if conf.VRSConnectionCheckTimer == 0 {
		log.Warnf("VRS Connection keep alive timer not set. Using default value")
		conf.VRSConnectionCheckTimer = 180
	}

	if conf.NuageSiteId == 0 {
		log.Warnf("SiteId not set. It will not be used when specifying metadata")
		conf.NuageSiteId = -1
	}
}

func IsVSPFunctional(vrsConnection vrsSdk.VRSConnection) bool {

	log.Debugf("Verifying VRS-VSC connection state")
	state, err := vrsConnection.GetControllerState()
	if err != nil {
		log.Errorf("failed getting controller state %v", err)
		return false
	}
	return state == vrsSdk.ControllerConnected
}
