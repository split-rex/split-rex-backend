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
	PaymentInfo types.PaymentInfo
}

// init -> random + going to be deleted
func (uf *UserFactory) Init(id uuid.UUID) {
	uf.ID = id
	
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
	uf.ID = uuid.MustParse("82405519-d6ca-45ce-b7d6-eeba1a66df59")

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
	id, _ := uuid.Parse("f6b22870-d9b8-43df-a302-2f2bfc035640")
	uf.ID = id
	uf.Name = "userA"
	uf.Email = "userA@gmail.com"
	uf.Username = "userA"
	uf.Password = types.EncryptedString("userA")
}

// userB (is on DB, not going to be deleted)
func (uf *UserFactory) UserB() {
	id, _ := uuid.Parse("421f9f47-f997-4d10-baaa-db8aacada674")
	uf.ID = id
	uf.Name = "userB"
	uf.Email = "userB@gmail.com"
	uf.Username = "userB"
	uf.Password = types.EncryptedString("userB")
}

// userC (is on DB, not going to be deleted)
func (uf *UserFactory) UserC() {
	id, _ := uuid.Parse("f2be1202-5b02-4e62-a05d-c2064bb4ba89")
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
