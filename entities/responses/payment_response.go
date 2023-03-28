package responses

import (
	"github.com/google/uuid"
)

type UnsettledPaymentResponse struct {
	PaymentID   uuid.UUID `json:"payment_id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"fullname"`
	Color       uint      `json:"color"`
	TotalUnpaid float64   `json:"total_unpaid"`
	TotalPaid   float64   `json:"total_paid"`
	Status      string    `json:"status"`
}

type UnconfirmedPaymentResponse struct {
	PaymentID uuid.UUID `json:"payment_id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"fullname"`
	Color     uint      `json:"color"`
	TotalPaid float64   `json:"total_paid"`
}
