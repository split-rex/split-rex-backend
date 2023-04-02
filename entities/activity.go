package entities

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ActivityID   uuid.UUID `gorm:"not null;unique"`
	ActivityType string    `gorm:"not null"`
	UserID       uuid.UUID `gorm:"not null"`
	Date         time.Time `gorm:"not null"`
	RedirectID   uuid.UUID `gorm:"not null"`
	DetailID     uuid.UUID `gorm:"not null"`
}

type GroupActivity struct {
	ActivityID uuid.UUID `gorm:"not null;unique"`
	GroupID    uuid.UUID `gorm:"not null"`
	UserID1    uuid.UUID `gorm:"not null"`
	UserID2    uuid.UUID `gorm:"not null"`
	Amount     float64   `gorm:"not null"`
	Date       time.Time `gorm:"not null"`
}

type PaymentActivity struct {
	PaymentActivityID uuid.UUID `gorm:"not null;unique"`
	Name              string    `gorm:"not null"`
	Status            string    `gorm:"not null"`
	Amount            float64   `gorm:"not null"`
	GroupName         string    `gorm:"not null"`
}

type TransactionActivity struct {
	TransactionActivityID uuid.UUID `gorm:"not null;unique"`
	Name                  string    `gorm:"not null"`
	GroupName             string    `gorm:"not null"`
}

type ReminderActivity struct {
	ReminderActivityID uuid.UUID `gorm:"not null;unique"`
	Name               string    `gorm:"not null"`
	GroupName          string    `gorm:"not null"`
}
