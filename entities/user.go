package entities

import (
	"split-rex-backend/types"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID             `gorm:"not null;unique"`
	Name     string                `gorm:"not null"`
	Email    string                `gorm:"not null;unique"`
	Username string                `gorm:"not null;unique"`
	Password types.EncryptedString `gorm:"not null"`
}
