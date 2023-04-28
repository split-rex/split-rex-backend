package requests

import (
	"time"

	"github.com/google/uuid"
)

type NotifRequest struct {
	UserID    uuid.UUID `json:"user_id" form:"user_id" query:"user_id"`
	GroupID   uuid.UUID `json:"group_id" form:"group_id" query:"group_id"`
	GroupName string    `json:"group_name" form:"group_name" query:"group_name"`
	Amount    float64   `json:"amount" form:"amount" query:"amount"`
	Name      string    `json:"name" form:"name" query:"name"`
	Date      time.Time `json:"date" form:"date" query:"date"`
}
