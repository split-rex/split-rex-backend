package responses

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	Notifications []NotificationDetail `json:"notifications"`
}

type NotificationDetail struct {
	NotificationID uuid.UUID `json:"notification_id"`
	GroupID        uuid.UUID `json:"group_id"`
	GroupName      string    `json:"group_name"`
	Amount         float64   `json:"amount"`
	Name           string    `json:"name"`
	Date           time.Time `json:"date"`
}
