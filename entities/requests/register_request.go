package requests

import (
	"split-rex-backend/types"
)

type RegisterRequest struct {
	Name     string                `gorm:"not null"`
	Email    string                `gorm:"not null;unique"`
	Username string                `gorm:"not null;unique"`
	Password types.EncryptedString `gorm:"not null"`
}
