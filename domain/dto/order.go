package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductUUID uuid.UUID `json:"productUUID" validate:"required"`
	Quantity    int       `json:"quantity" validate:"required,gt=0"`
}

type OrderItemResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	ProductID   uint      `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Subtotal    float64   `json:"subtotal"`
}

type OrderResponse struct {
	UUID      uuid.UUID           `json:"uuid"`
	UserID    uint                `json:"userId"`
	Amount    float64             `json:"amount"`
	Status    string              `json:"status"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt *time.Time          `json:"createdAt,omitempty"`
	PaidAt    *time.Time          `json:"paidAt,omitempty"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
}
