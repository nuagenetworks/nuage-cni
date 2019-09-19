package table

import (
	"github.com/nuagenetworks/libvrsdk/ovsdb"
	"github.com/socketplane/libovsdb"
)

//FakeTable is a fake table that can be used for testing
type FakeTable struct {
	ReadRowsFunc func(ovs *libovsdb.OvsdbClient, readRowArgs ovsdb.ReadRowArgs) ([]map[string]interface{}, error)
}

//InsertRow inserts a row into table
func (f *FakeTable) InsertRow(ovs *libovsdb.OvsdbClient, row ovsdb.NuageTableRow) error {
	return nil
}

//DeleteRow deletes a row from table
func (f *FakeTable) DeleteRow(ovs *libovsdb.OvsdbClient, condition []string) error {
	return nil
}

//ReadRow read one row matching readRowArgs
func (f *FakeTable) ReadRow(ovs *libovsdb.OvsdbClient, readRowArgs ovsdb.ReadRowArgs) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	return m, nil
}

//ReadRows returns all the rows matching readRowArgs
func (f *FakeTable) ReadRows(ovs *libovsdb.OvsdbClient, readRowArgs ovsdb.ReadRowArgs) ([]map[string]interface{}, error) {
	if f.ReadRowsFunc != nil {
		return f.ReadRowsFunc(ovs, readRowArgs)
	}
	m := make([]map[string]interface{}, 10)
	return m, nil
}

//UpdateRow updates the row in the table
func (f *FakeTable) UpdateRow(ovs *libovsdb.OvsdbClient, ovsdbRow map[string]interface{}, condition []string) error {
	return nil
}
