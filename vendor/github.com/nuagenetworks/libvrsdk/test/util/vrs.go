package util

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

const (
	add    = "add"
	delete = "del"
)

// EnableOVSDBRPCSocket will add an interface to the ovsdb-server
// to make it accept RPCs via TCP socket
func EnableOVSDBRPCSocket(port int) error {

	cmdstr := fmt.Sprintf("ovs-appctl -t ovsdb-server ovsdb-server/add-remote ptcp:%d", port)
	cmd := exec.Command("bash", "-c", cmdstr)
	_, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("Error while add an interface to the ovsdb-server to make it accept RPCs via TCP socket%v", err)
	}

	return nil
}

// CreateVETHPair will help user create veth pairs to associate
// with a VM or a Container
func CreateVETHPair(portList []string) error {

	cmdstr := fmt.Sprintf("ip link %s %s type veth peer name %s", add, portList[0], portList[1])
	cmd := exec.Command("bash", "-c", cmdstr)
	_, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("Error while creating veth pair on VRS %v", err)
	}

	for index := range portList {
		cmdstr = fmt.Sprintf("ip link set dev %s up", portList[index])
		cmd = exec.Command("bash", "-c", cmdstr)
		_, err = cmd.Output()

		if err != nil {
			return fmt.Errorf("Error while bringing up veth interface on VRS %v", err)
		}
	}

	return nil
}

// DeleteVETHPair will help user delete veth pairs on VRS
func DeleteVETHPair(entityPort string, brPort string) error {

	cmdstr := fmt.Sprintf("ip link %s %s type veth peer name %s", delete, entityPort, brPort)
	cmd := exec.Command("bash", "-c", cmdstr)
	_, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("Error while creating veth pair on VRS %v", err)
	}

	return nil
}

// AddVETHPortToVRS will help add veth ports to VRS alubr0
func AddVETHPortToVRS(port string, vmuuid string, vmname string) error {

	cmdstr := fmt.Sprintf("/usr/bin/ovs-vsctl --no-wait --if-exists del-port alubr0 %s -- %s-port alubr0 %s -- set interface %s 'external-ids={vm-uuid=%s,vm-name=%s}'", port, add, port, port, vmuuid, vmname)
	cmd := exec.Command("bash", "-c", cmdstr)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Problem adding veth port to alubr0 on VRS %v", err)
	}

	return nil
}

// RemoveVETHPortFromVRS will help delete veth ports from VRS alubr0
func RemoveVETHPortFromVRS(port string) error {

	cmdstr := fmt.Sprintf("/usr/bin/ovs-vsctl --no-wait %s-port alubr0 %s", delete, port)
	cmd := exec.Command("bash", "-c", cmdstr)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Problem deleting veth port from alubr0 on VRS %v", err)
	}

	return nil
}

// GenerateMAC will act as a pseudo random MAC generator
func GenerateMAC() string {

	arr := make([]int, 6)
	var num int
	for i := 0; i < 6; i++ {
		for {
			num = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
			if num >= 10 && num <= 99 {
				break
			}
		}

		arr[i] = num
	}

	mac := fmt.Sprintf("%d:%d:%d:%d:%d:%d", arr[0], arr[1], arr[2], arr[3], arr[4], arr[5])
	return mac
}
