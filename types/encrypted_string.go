package types

import (
	"database/sql/driver"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type EncryptedString []byte

func (es *EncryptedString) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("")
	}

	*es = bytes
	return nil
}

func (es EncryptedString) Value() (driver.Value, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(es, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func (EncryptedString) GormDataType() string {
	return "bytea"
}
