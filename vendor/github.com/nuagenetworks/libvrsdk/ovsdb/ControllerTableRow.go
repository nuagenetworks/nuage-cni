package ovsdb

const (
	//ControllerTable is the table name
	ControllerTable = "Controller"
	//ControllerTableColumnRole column role in controller table
	ControllerTableColumnRole = "role"
)

//ControllerTableRow represents a row in Controller Table
type ControllerTableRow struct {
	Role string
}

//Equals checks for equality of two rows in Controller Table
func (row *ControllerTableRow) Equals(otherRow interface{}) bool {
	controllerTableRow, ok := otherRow.(ControllerTableRow)
	if !ok {
		return false
	}

	if row.Role != controllerTableRow.Role {
		return false
	}
	return true
}
