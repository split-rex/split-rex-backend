package entities

import (
	"split-rex-backend/types"

	"github.com/google/uuid"
)

type Friend struct {
	ID           uuid.UUID `gorm:"not null;unique"`
	Friend_id    types.ArrayOfUUID
	Req_received types.ArrayOfUUID
	Req_sent     types.ArrayOfUUID
}
