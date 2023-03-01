package entities

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	GroupID   uuid.UUID   `gorm:"not null;unique"`
	Name      string      `gorm:"not null"`
	MemberID  []uuid.UUID `gorm:"not null;type:"`
	StartDate time.Time   `gorm:"not null"`
	EndDate   time.Time   `gorm:"not null"`
}
