package requests

import (
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
)

type UserCreateTransactionRequest struct {
	Name        string            `json:"name" form:"name" query:"name"`
	Description string            `json:"description" form:"description" query:"description"`
	GroupID     uuid.UUID         `json:"group_id" form:"group_id" query:"group_id"`
	Date        time.Time         `json:"date" form:"date" query:"date" default:"time.Now()"`
	Subtotal    float64           `json:"subtotal" form:"subtotal" query:"subtotal" default:"0.0"`
	Tax         float64           `json:"tax" form:"tax" query:"tax" default:"0.0"`
	Service     float64           `json:"service" form:"service" query:"service" default:"0.0"`
	Total       float64           `json:"total" form:"total" query:"total" default:"0.0"`
	BillOwner   uuid.UUID         `json:"bill_owner" form:"bill_owner" query:"bill_owner"`
	Items       types.ArrayOfUUID `json:"items" form:"items" query:"items"`
}

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
