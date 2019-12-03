package api

import (
	"fmt"

	"github.com/nuagenetworks/libvrsdk/ovsdb"
)

//ControllerState is for connection state
type ControllerState string

const (
	//ControllerConnected if ovs is connected to controller
	ControllerConnected ControllerState = "connected"
	//ControllerDisconnected if ovs is disconnected from controller
	ControllerDisconnected ControllerState = "disconnected"
	//ControllerStateUnknown if ovs is disconnected from controller
	ControllerStateUnknown ControllerState = "unknown"
	//MasterController is the master controller for this ovs
	MasterController string = "master"
)

//GetControllerState return the state of the controller connection
func (vrsConnection *VRSConnection) GetControllerState() (ControllerState, error) {

	readRowArgs := ovsdb.ReadRowArgs{
		Columns:   []string{ovsdb.ControllerTableColumnRole},
		Condition: []string{ovsdb.ControllerTableColumnRole, "==", MasterController},
	}

	rows, err := vrsConnection.controllerTable.ReadRows(vrsConnection.ovsdbClient, readRowArgs)
	if err != nil {
		return ControllerStateUnknown, fmt.Errorf("Unable to controller state info %v", err)
	}

	if len(rows) != 1 {
		return ControllerDisconnected, nil
	}

	return ControllerConnected, nil
}
