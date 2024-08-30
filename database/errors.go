package database

import "fmt"

type NoID struct {
	objectName string
}

func (e NoID) Error() string {
	return fmt.Sprintf("ERROR: '%v' has no id yet! Maybe new or not in the db\n", e.objectName)
}
