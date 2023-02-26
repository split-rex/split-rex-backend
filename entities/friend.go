package entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Friend struct {
	ID           uuid.UUID      `gorm:"not null;unique"`
	Friend_id    pq.StringArray `gorm:"type:text[]"`
	Req_received pq.StringArray `gorm:"type:text[]"`
	Req_sent     pq.StringArray `gorm:"type:text[]"`
}
