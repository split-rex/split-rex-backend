package responses

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDetailResponse struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	GroupID     uuid.UUID            `json:"group_id"`
	GroupName   string               `json:"group_name"`
	Date        time.Time            `json:"date"`
	Items       []ItemDetailResponse `json:"items"`
	Subtotal    float64              `json:"subtotal"`
	Tax         float64              `json:"tax"`
	Service     float64              `json:"service"`
	Total       float64              `json:"total"`
}

type ItemDetailResponse struct {
	ItemID     uuid.UUID                `json:"item_id"`
	Name       string                   `json:"name"`
	Quantity   int                      `json:"quantity"`
	Price      float64                  `json:"price"`
	TotalPrice float64                  `json:"total_price"`
	Consumer   []ConsumerDetailResponse `json:"consumer"`
}

type ConsumerDetailResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Color  uint      `json:"color"`
}
