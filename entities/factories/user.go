package factories

import (
	"split-rex-backend/types"

	"github.com/google/uuid"
)

type UserFactory struct {
	ID       uuid.UUID             
	Name     string                `default:"testing"`
	Email    string                `default:"testing@gmail.com"`
	Username string                `default:"testing"`
	Color    uint                  `default:"1"`
	Password types.EncryptedString `password:"testing"`
	Groups   types.ArrayOfUUID
}
