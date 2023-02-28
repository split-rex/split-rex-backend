package types

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type ArrayOfUUID []uuid.UUID

func (arrayOfUUID *ArrayOfUUID) Scan(value interface{}) error {
	*arrayOfUUID = []uuid.UUID{}
	for _, uid := range value.([]uuid.UUID) {
		*arrayOfUUID = append(*arrayOfUUID, uid)
	}

	return nil
}

func (arrayOfUUID ArrayOfUUID) Value() (driver.Value, error) {
	var values []uuid.UUID
	for _, uid := range arrayOfUUID {
		values = append(values, uid)
	}
	return values, nil
}

func (ArrayOfUUID) GormDataType() string {
	return "uuid[]"
}
