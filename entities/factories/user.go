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

// userAuth (is on DB, not going to be deleted), this user is for updating endpoint
func (uf *UserFactory) InitAuth() {
	uf.ID = uuid.MustParse("24ba7892-ea12-4d20-92ca-02e016ee711a")

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
	id, _ := uuid.Parse("9368048d-7e77-4fc1-b81e-a5abfef844e1")
	uf.ID = id
	uf.Name = "userA"
	uf.Email = "userA@gmail.com"
	uf.Username = "userA"
	uf.Password = types.EncryptedString("userA")
}

// userB (is on DB, not going to be deleted)
func (uf *UserFactory) UserB() {
	id, _ := uuid.Parse("3af1e0b7-2a57-4834-b8ba-a7d5f3f5da8b")
	uf.ID = id
	uf.Name = "userB"
	uf.Email = "userB@gmail.com"
	uf.Username = "userB"
	uf.Password = types.EncryptedString("userB")
}

// userC (is on DB, not going to be deleted)
func (uf *UserFactory) UserC() {
	id, _ := uuid.Parse("af902382-60d8-4dc6-bd76-dc3f1f061e7a")
	uf.ID = id
	uf.Name = "userC"
	uf.Email = "userC@gmail.com"
	uf.Username = "userC"
	uf.Password = types.EncryptedString("userC")
}

// userD (is on DB, not going to be deleted)
func (uf *UserFactory) UserD() {
	id, _ := uuid.Parse("e7d56a3d-930f-45aa-9fcf-a154f2e2db8c")
	uf.ID = id
	uf.Name = "userD"
	uf.Email = "userD@gmail.com"
	uf.Username = "userD"
	uf.Password = types.EncryptedString("userD")
}
