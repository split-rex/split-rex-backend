package entities

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ExpenseID uuid.UUID `gorm:"primaryKey;not null;unique"`
	UserID    uuid.UUID `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
}
