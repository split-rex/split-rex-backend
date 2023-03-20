package responses

import (
	"github.com/google/uuid"
)

type UnsettledTransactionResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"fullname"`
	TotalUnpaid float64   `json:"total_unpaid"`
	TotalPaid   float64   `json:"total_paid"`
	Status      string    `json:"status"`
}
