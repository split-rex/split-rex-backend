package factories

import (
	"split-rex-backend/types"

	"github.com/google/uuid"
)

type UserFactory struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Username string
	Color    uint
	Password types.EncryptedString
	Groups   types.ArrayOfUUID
}

func (uf *UserFactory) Init() {

	if uf.Name == "" {
		uf.Name = "ABC"
	}

	if uf.Email == "" {
		uf.Email = "testing@gmail.com"
	}

	if uf.Username == "" {
		uf.Username = "testing"
	}

	if uf.Password == nil {
		uf.Password = types.EncryptedString("testing")
	}
}
