package dto

import (
	"time"

	"github.com/google/uuid"
)

type PaymentRequest struct {
}

type PaymentResponse struct {
	UUID          uuid.UUID  `json:"uuid"`
	OrderID       uint       `json:"orderId"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"`
	TransactionID *string    `json:"transactionId,omitempty"`
	PaidAt        *time.Time `json:"paidAt,omitempty"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
}
