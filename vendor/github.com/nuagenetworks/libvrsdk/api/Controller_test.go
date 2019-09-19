package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nuagenetworks/libvrsdk/ovsdb"
	testutils "github.com/nuagenetworks/libvrsdk/test/util/table"
	"github.com/socketplane/libovsdb"
)

type controllerTableTestVec struct {
	readRows func(ovs *libovsdb.OvsdbClient, readRowArgs ovsdb.ReadRowArgs) ([]map[string]interface{}, error)
	out      ControllerState
	err      error
}

func TestGetControllerState(t *testing.T) {

	fakeTable := &testutils.FakeTable{}

	vrsConnection := VRSConnection{
		controllerTable: fakeTable,
	}

	vec := []controllerTableTestVec{
		{func(ovs *libovsdb.OvsdbClient, readRowArgs ovsdb.ReadRowArgs) ([]map[string]interface{}, error) {
			m := make([]map[string]interface{}, 1)
			return m, nil
		}, ControllerConnected, nil},
		{func(ovs *libovsdb.OvsdbClient, readRowArgs ovsdb.ReadRowArgs) ([]map[string]interface{}, error) {
			m := make([]map[string]interface{}, 0)
			return m, nil
		}, ControllerDisconnected, nil},
		{func(ovs *libovsdb.OvsdbClient, readRowArgs ovsdb.ReadRowArgs) ([]map[string]interface{}, error) {
			return nil, fmt.Errorf("some random error")
		}, ControllerStateUnknown, fmt.Errorf("some random error")},
	}

	for _, tt := range vec {
		fakeTable.ReadRowsFunc = tt.readRows
		state, err := vrsConnection.GetControllerState()
		if tt.out != state {
			t.Fatalf("controller state %v did not match expected %v", state, tt.out)
		}
		if err != nil && tt.err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("errors did not match: expected %v found %v", err.Error(), tt.err.Error())
		}
		if tt.err != nil && err == nil {
			t.Fatalf("expected non nil error %v but found", tt.err)
		}
		if tt.err == nil && err != nil {
			t.Fatalf("expected nil error but found %v", err)
		}

	}
}
