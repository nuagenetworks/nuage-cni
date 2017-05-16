package api

import (
	"errors"
	"fmt"
	"github.com/nuagenetworks/libvrsdk/api/port"
	"github.com/nuagenetworks/libvrsdk/ovsdb"
	"github.com/socketplane/libovsdb"
	"reflect"
)

type empty struct{}

// PortIPv4Info defines details to be populated
// for container port resolved in OVSDB
type PortIPv4Info struct {
	IPAddr     string
	Gateway    string
	Mask       string
	MAC        string
	Registered bool
}

// GetAllPorts returns the slice of all the vport names attached to the VRS
func (vrsConnection *VRSConnection) GetAllPorts() ([]string, error) {

	readRowArgs := ovsdb.ReadRowArgs{
		Condition: []string{ovsdb.NuagePortTableColumnName, "!=", "xxxx"},
		Columns:   []string{ovsdb.NuagePortTableColumnName},
	}

	var nameRows []map[string]interface{}
	var err error
	if nameRows, err = vrsConnection.portTable.ReadRows(vrsConnection.ovsdbClient, readRowArgs); err != nil {
		return nil, fmt.Errorf("Unable to obtain the entity names %v", err)
	}

	var names []string
	for _, name := range nameRows {
		names = append(names, name[ovsdb.NuagePortTableColumnName].(string))
	}

	return names, nil
}

// CreatePort creates a new vPort in the Nuage VRS. The only mandatory inputs required to create
// a port are it's name and MAC address
func (vrsConnection *VRSConnection) CreatePort(name string, attributes port.Attributes,
	metadata map[port.MetadataKey]string) error {

	portMetadata := make(map[string]string)

	for k, v := range metadata {
		portMetadata[string(k)] = v
	}

	nuagePortRow := ovsdb.NuagePortTableRow{
		Name:             name,
		Mac:              attributes.MAC,
		Bridge:           attributes.Bridge,
		NuageDomain:      metadata[port.MetadataKeyDomain],
		NuageNetwork:     metadata[port.MetadataKeyNetwork],
		NuageNetworkType: metadata[port.MetadataKeyNetworkType],
		NuageZone:        metadata[port.MetadataKeyZone],
		VMDomain:         attributes.Platform,
		Metadata:         portMetadata,
	}

	if err := vrsConnection.portTable.InsertRow(vrsConnection.ovsdbClient, &nuagePortRow); err != nil {
		return fmt.Errorf("Problem adding port info to VRS %v", err)
	}

	return nil
}

// DestroyPort purges a port from the Nuage VRS
func (vrsConnection *VRSConnection) DestroyPort(name string) error {

	condition := []string{ovsdb.NuagePortTableColumnName, "==", name}
	if err := vrsConnection.portTable.DeleteRow(vrsConnection.ovsdbClient, condition); err != nil {
		return fmt.Errorf("Unable to remove the port from VRS %v", err)
	}

	return nil
}

// GetPortState gets the current resolution state of the port namely the IP address, Subnet Mask, Gateway,
// EVPN ID and VRF ID
func (vrsConnection VRSConnection) GetPortState(name string) (map[port.StateKey]interface{}, error) {

	readRowArgs := ovsdb.ReadRowArgs{
		Columns: []string{ovsdb.NuagePortTableColumnIPAddress, ovsdb.NuagePortTableColumnSubnetMask,
			ovsdb.NuagePortTableColumnGateway, ovsdb.NuagePortTableColumnEVPNID,
			ovsdb.NuagePortTableColumnVRFId},
		Condition: []string{ovsdb.NuagePortTableColumnName, "==", name},
	}

	var row map[string]interface{}
	var err error
	if row, err = vrsConnection.portTable.ReadRow(vrsConnection.ovsdbClient, readRowArgs); err != nil {
		return make(map[port.StateKey]interface{}), fmt.Errorf("Unable to obtain the port row %v", err)
	}

	portState := make(map[port.StateKey]interface{})
	portState[port.StateKeyIPAddress] = row[ovsdb.NuagePortTableColumnIPAddress]
	portState[port.StateKeySubnetMask] = row[ovsdb.NuagePortTableColumnSubnetMask]
	portState[port.StateKeyGateway] = row[ovsdb.NuagePortTableColumnGateway]
	portState[port.StateKeyVrfID] = row[ovsdb.NuagePortTableColumnVRFId]
	portState[port.StateKeyEvpnID] = row[ovsdb.NuagePortTableColumnEVPNID]

	return portState, nil
}

// UpdatePortAttributes updates the attributes of the vPort
func (vrsConnection *VRSConnection) UpdatePortAttributes(name string, attrs port.Attributes) error {
	row := make(map[string]interface{})

	row[ovsdb.NuagePortTableColumnBridge] = attrs.Bridge
	row[ovsdb.NuagePortTableColumnMAC] = attrs.MAC
	row[ovsdb.NuagePortTableColumnVMDomain] = attrs.Platform

	condition := []string{ovsdb.NuagePortTableColumnName, "==", name}

	if err := vrsConnection.portTable.UpdateRow(vrsConnection.ovsdbClient, row, condition); err != nil {
		return fmt.Errorf("Unable to update the port attributes %s %v %v", name, attrs, err)
	}

	return nil
}

// UpdatePortMetadata updates the metadata for the vPort
func (vrsConnection *VRSConnection) UpdatePortMetadata(name string, metadata map[string]string) error {
	row := make(map[string]interface{})

	metadataOVSDB, err := libovsdb.NewOvsMap(metadata)
	if err != nil {
		return fmt.Errorf("Unable to create OVSDB map %v", err)
	}

	row[ovsdb.NuagePortTableColumnMetadata] = metadataOVSDB

	key := string(port.MetadataKeyDomain)
	if len(metadata[key]) != 0 {
		row[ovsdb.NuagePortTableColumnNuageDomain] = metadata[key]
		delete(metadata, key)
	}

	key = string(port.MetadataKeyNetwork)
	if len(metadata[key]) != 0 {
		row[ovsdb.NuagePortTableColumnNuageNetwork] = metadata[key]
		delete(metadata, key)
	}

	key = string(port.MetadataKeyZone)
	if len(metadata[key]) != 0 {
		row[ovsdb.NuagePortTableColumnNuageZone] = metadata[key]
		delete(metadata, key)
	}

	condition := []string{ovsdb.NuagePortTableColumnName, "==", name}

	if err := vrsConnection.portTable.UpdateRow(vrsConnection.ovsdbClient, row, condition); err != nil {
		return fmt.Errorf("Unable to update the port metadata %s %v %v", name, metadata, err)
	}

	return nil
}

// RegisterForPortUpdates will help register via channel
// for VRS port table updates
func (vrsConnection *VRSConnection) RegisterForPortUpdates(brport string, pnc chan *PortIPv4Info) error {
	vrsConnection.registrationChannel <- &Registration{Brport: brport, Channel: pnc, Register: true}
	return nil
}

// DeregisterForPortUpdates will help de-register for VRS port table updates
func (vrsConnection *VRSConnection) DeregisterForPortUpdates(brport string) error {
	vrsConnection.registrationChannel <- &Registration{Brport: brport, Channel: nil, Register: false}
	return nil
}

func (vrsConnection VRSConnection) handlePortRegistration(registration *Registration) error {
	brport := registration.Brport
	register := registration.Register
	pnc := registration.Channel
	if register {
		if _, ok := vrsConnection.pncTable[brport]; ok {
			return fmt.Errorf("Already registered for this bridge port %s", brport)
		}
		vrsConnection.pncTable[brport] = pnc
		if portInfo, exists := vrsConnection.pnpTable[brport]; exists {
			select {
			case pnc <- &portInfo:
			default:
			}
			delete(vrsConnection.pnpTable, brport)
		}
	} else {
		delete(vrsConnection.pncTable, brport)
	}
	return nil
}

func (vrsConnection VRSConnection) getPortInfo(row *libovsdb.Row) (*PortIPv4Info, error) {
	portIPv4Info := PortIPv4Info{Registered: true}
	if _, ok := row.Fields["ip_addr"]; ok {
		ip := row.Fields["ip_addr"].(string)
		if ip != "" {
			portIPv4Info.IPAddr = ip
		} else {
			return nil, errors.New("Invalid or empty ip")
		}
	}
	if _, ok := row.Fields["subnet_mask"]; ok {
		subnet := row.Fields["subnet_mask"].(string)
		if subnet != "" {
			portIPv4Info.Mask = subnet
		} else {
			return nil, errors.New("Invalid or empty subnet")
		}
	}
	if _, ok := row.Fields["gateway"]; ok {
		gw := row.Fields["gateway"].(string)
		if gw != "" {
			portIPv4Info.Gateway = gw
		} else {
			return nil, errors.New("Invalid or empty gateway")
		}
	}
	if _, ok := row.Fields["mac"]; ok {
		mac := row.Fields["mac"].(string)
		if mac != "" {
			portIPv4Info.MAC = mac
		} else {
			return nil, errors.New("Invalid or empty port MAC address")
		}
	}

	return &portIPv4Info, nil
}

func (vrsConnection VRSConnection) processUpdates(updates *libovsdb.TableUpdates) error {
	for _, tableUpdate := range updates.Updates {
		for _, row := range tableUpdate.Rows {
			empty := libovsdb.Row{}
			if !reflect.DeepEqual(row.New, empty) {
				//check for whether the port is already registered for updates
				portInfo, err := vrsConnection.getPortInfo(&(row.New))
				if err == nil {
					if _, ok := (row.New).Fields["name"]; ok {
						portName := (row.New).Fields["name"].(string)
						if pncChannel, exists := vrsConnection.pncTable[portName]; exists {
							select {
							case pncChannel <- portInfo:
							default:
							}
						} else {
							vrsConnection.pnpTable[portName] = *portInfo
						}
					}
				}
			} else { //delete case
				if _, ok := (row.Old).Fields["name"]; ok {
					portName := (row.Old).Fields["name"].(string)
					if pncChannel, exists := vrsConnection.pncTable[portName]; exists {
						select {
						case pncChannel <- &PortIPv4Info{Registered: false}:
						default:
						}
						delete(vrsConnection.pncTable, portName)
					}
					delete(vrsConnection.pnpTable, portName)
				}
			}
		}
	}
	return nil
}
