package entities

import (
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	TransactionID uuid.UUID         `gorm:"not null;unique"`
	Name          string            `gorm:"not null"`
	Description   string            `gorm:"not null"`
	GroupID       uuid.UUID         `gorm:"not null"`
	Date          time.Time         `gorm:"not null"`
	Subtotal      float64           `gorm:"not null"`
	Tax           float64           `gorm:"not null"`
	Service       float64           `gorm:"not null"`
	Total         float64           `gorm:"not null"`
	BillOwner     uuid.UUID         `gorm:"not null"`
	Items         types.ArrayOfUUID `gorm:"not null"`
}
