package factories

import (
	"split-rex-backend/types"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
)

type UserFactory struct {
	ID          uuid.UUID
	Name        string
	Email       string
	Username    string
	Color       uint
	Password    types.EncryptedString
	Groups      types.ArrayOfUUID
	PaymentInfo map[string][]map[int]string
}

// init -> random + going to be deleted
func (uf *UserFactory) Init(id uuid.UUID) {
	uf.ID = uuid.New()

	if uf.Name == "" {
		uf.Name = faker.Name()
	}

	if uf.Email == "" {
		uf.Email = faker.Email()
	}

	if uf.Username == "" {
		uf.Username = faker.Username()
	}

	if uf.Password == nil {
		uf.Password = types.EncryptedString(faker.Password())
	}
}

// userAuth (is on DB, not going to be deleted)
func (uf *UserFactory) InitAuth() {

	if uf.Name == "" {
		uf.Name = "auth_testing"
	}

	if uf.Email == "" {
		uf.Email = "auth_testing@gmail.com"
	}

	if uf.Username == "" {
		uf.Username = "auth_testing"
	}

	if uf.Password == nil {
		uf.Password = types.EncryptedString("auth_testing")
	}
}

// userA (is on DB, not going to be deleted)
func (uf *UserFactory) UserA() {
	id, _ := uuid.Parse("cf734de2-2952-4766-88f9-bfae95e1c2f0")
	uf.ID = id
	uf.Name = "userA"
	uf.Email = "userA@gmail.com"
	uf.Username = "userA"
	uf.Password = types.EncryptedString("userA")
}

// userB (is on DB, not going to be deleted)
func (uf *UserFactory) UserB() {
	id, _ := uuid.Parse("06c2e522-30e9-4171-8efb-9d27b7c4bee9")
	uf.ID = id
	uf.Name = "userB"
	uf.Email = "userB@gmail.com"
	uf.Username = "userB"
	uf.Password = types.EncryptedString("userB")
}

// userC (is on DB, not going to be deleted)
func (uf *UserFactory) UserC() {
	id, _ := uuid.Parse("acbe5a63-1390-41e1-b463-7c9b2b2a0f46")
	uf.ID = id
	uf.Name = "userC"
	uf.Email = "userC@gmail.com"
	uf.Username = "userC"
	uf.Password = types.EncryptedString("userC")
}

// userD (is on DB, not going to be deleted)
func (uf *UserFactory) UserD() {
	id, _ := uuid.Parse("cf4dda6a-a3b6-47c8-b7a9-035e43b4967a")
	uf.ID = id
	uf.Name = "userD"
	uf.Email = "userD@gmail.com"
	uf.Username = "userD"
	uf.Password = types.EncryptedString("userD")
}
