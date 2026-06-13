package model

import "time"

const (
	PaymentStatusSuccess = "SUCCESS"
	PaymentStatusFailed  = "FAILED"
)

type Payment struct {
	ID        int64     `json:"id"`
	PaymentNo string    `json:"payment_no"`
	OrderID   int64     `json:"order_id"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
