package entities

import (
	"github.com/google/uuid"
)

type Friend struct {
	ID           uuid.UUID   `gorm:"not null;unique"`
	Friend_id    []uuid.UUID `gorm:"type:bytea"`
	Req_received []uuid.UUID `gorm:"type:bytea"`
	Req_sent     []uuid.UUID `gorm:"type:bytea"`
}
