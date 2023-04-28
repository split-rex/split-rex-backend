package entities

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	NotificationID uuid.UUID `gorm:"primaryKey"`
	GroupID        uuid.UUID `gorm:"not null"`
	GroupName      string    `gorm:"not null"`
	Amount         float64   `gorm:"not null"`
	Name           string    `gorm:"not null"`
	Date           time.Time `gorm:"not null"`
}
