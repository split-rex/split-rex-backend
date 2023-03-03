package entities

import (
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	GroupID   uuid.UUID         `gorm:"primaryKey"`
	Name      string            `gorm:"not null"`
	MemberID  types.ArrayOfUUID `gorm:"not null"`
	StartDate time.Time         `gorm:"not null"`
	EndDate   time.Time         `gorm:"not null"`
}
