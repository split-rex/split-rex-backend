package requests

import (
	"github.com/google/uuid"
)

type UpdatePaymentRequest struct {
	GroupID     uuid.UUID        `json:"group_id" form:"group_id" query:"group_id"`
	ListPayment []PaymentRequest `json:"list_payment" form:"list_payment" query:"list_payment"`
}

type PaymentRequest struct {
	UserID      uuid.UUID `json:"user_id" form:"user_id" query:"user_id"`
	TotalUnpaid float64   `json:"total_unpaid" form:"total_unpaid" query:"total_unpaid"`
}

type ResolveTransactionRequest struct {
	GroupID uuid.UUID `json:"group_id" form:"group_id" query:"group_id"`
}

type SettleRequest struct {
	PaymentID uuid.UUID `json:"payment_id" form:"payment_id" query:"payment_id"`
	TotalPaid float64   `json:"total_paid" form:"total_paid" query:"total_paid"`
}

type ConfirmRequest struct {
	PaymentID uuid.UUID `json:"payment_id" form:"payment_id" query:"payment_id"`
}
