package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

type ArrayOfUUID []uuid.UUID

func (arrayOfUUID ArrayOfUUID) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &arrayOfUUID)
}

func (arrayOfUUID ArrayOfUUID) Value() (driver.Value, error) {
	val, err := json.Marshal(arrayOfUUID)
	return string(val), err
}

func (ArrayOfUUID) GormDataType() string {
	return "string"
}

func (arrayOfUUID ArrayOfUUID) Count() int {
	return arrayOfUUID.Count()
}

func (arrayOfUUID ArrayOfUUID) Contains(id uuid.UUID) bool {
	for _, uid := range arrayOfUUID {
		if uid == id {
			return true
		}
	}
	return false
}
