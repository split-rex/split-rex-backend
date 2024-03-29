package requests

import (
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
)

type UserCreateTransactionRequest struct {
	Name        string        `json:"name" form:"name" query:"name"`
	Description string        `json:"description" form:"description" query:"description"`
	GroupID     uuid.UUID     `json:"group_id" form:"group_id" query:"group_id"`
	Date        time.Time     `json:"date" form:"date" query:"date" default:"time.Now()"`
	Subtotal    float64       `json:"subtotal" form:"subtotal" query:"subtotal" default:"0.0"`
	Tax         float64       `json:"tax" form:"tax" query:"tax" default:"0.0"`
	Service     float64       `json:"service" form:"service" query:"service" default:"0.0"`
	Total       float64       `json:"total" form:"total" query:"total" default:"0.0"`
	BillOwner   uuid.UUID     `json:"bill_owner" form:"bill_owner" query:"bill_owner"`
	Items       []ItemRequest `json:"items" form:"items" query:"items"`
}

type ItemRequest struct {
	Name     string            `json:"name" form:"name" query:"name"`
	Quantity int               `json:"quantity" form:"quantity" query:"quantity"`
	Price    float64           `json:"price" form:"price" query:"price"`
	Consumer types.ArrayOfUUID `json:"consumer" form:"consumer" query:"consumer"`
}
