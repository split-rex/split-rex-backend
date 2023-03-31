package responses

import (
	"time"

	"github.com/google/uuid"
)

type ActivityResponse[T interface{}] struct {
	ActivityID   uuid.UUID `json:"activity_id"`
	ActivityType string    `json:"activity_type"`
	Date         time.Time `json:"date"`
	RedirectID   uuid.UUID `json:"redirect_id"`
	Detail       T         `json:"detail"`
}

type GroupActivityResponse struct {
	ActivityID uuid.UUID `json:"activity_id"`
	Date       time.Time `json:"date"`
	Name1      string    `json:"name1"`
	Name2      string    `json:"name2"`
	Amount     float64   `json:"amount"`
}

type PaymentActivityResponse struct {
	PaymentActivityID uuid.UUID `json:"payment_activity_id"`
	Name              string    `json:"name"`
	Status            string    `json:"status"`
	Amount            float64   `json:"amount"`
	GroupName         string    `json:"group_name"`
}

type TransactionActivityResponse struct {
	TransactionActivityID uuid.UUID `json:"transaction_activity_id"`
	Name                  string    `json:"name"`
	GroupName             string    `json:"group_name"`
}

type ReminderActivityResponse struct {
	ReminderActivityID uuid.UUID `json:"reminder_activity_id"`
	Name               string    `json:"name"`
	GroupName          string    `json:"group_name"`
}
