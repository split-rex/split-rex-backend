package entities

import (
	"split-rex-backend/types"

	"github.com/google/uuid"
)

type Item struct {
	ItemID   uuid.UUID         `gorm:"not null;unique"`
	Name     string            `gorm:"not null"`
	Quantity int               `gorm:"not null"`
	Price    float64           `gorm:"not null"`
	Consumer types.ArrayOfUUID `gorm:"not null"`
}
