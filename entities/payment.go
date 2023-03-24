package entities

import (
	"github.com/google/uuid"
)

type Payment struct {
	PaymentID   uuid.UUID `gorm:"primaryKey;not null;unique"`
	GroupID     uuid.UUID `gorm:"not null"`
	UserID1     uuid.UUID `gorm:"not null"`
	UserID2     uuid.UUID `gorm:"not null"`
	TotalUnpaid float64   `gorm:"not null"`
	TotalPaid   float64   `gorm:"not null"`
	Status      string    `gorm:"not null"`
}
