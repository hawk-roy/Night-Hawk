package model

import "time"

const (
	OrderStatusPendingPayment = "PENDING_PAYMENT"
)

type Order struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Username    string    `json:"username"`
	ProductID   int64     `json:"product_id"`
	ProductName string    `json:"product_name"`
	UnitPrice   int64     `json:"unit_price"`
	Quantity    int64     `json:"quantity"`
	TotalAmount int64     `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
