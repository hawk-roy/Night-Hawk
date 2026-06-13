package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hawk-roy/Night-Hawk/internal/model"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderNotPendingPayment = errors.New("order is not pending payment")
var ErrInvalidPaymentResult = errors.New("invalid payment result")

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func generatePaymentNo() string {
	return fmt.Sprintf("PAY%d", time.Now().UnixNano())
}

func (r *PaymentRepository) MockPayOrder(ctx context.Context, userID int64, orderID int64, result string) (*model.Payment, string, error) {
	if result != model.PaymentStatusSuccess && result != model.PaymentStatusFailed {
		return nil, "", ErrInvalidPaymentResult
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()

	var order model.Order
	row := tx.QueryRowContext(ctx, `
		SELECT id, order_no, user_id, total_amount, status
		FROM orders
		WHERE id = ?
		FOR UPDATE
	`, orderID)

	if err := row.Scan(&order.ID, &order.OrderNo, &order.UserID, &order.TotalAmount, &order.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrOrderNotFound
		}
		return nil, "", err
	}

	if order.UserID != userID {
		return nil, "", ErrOrderNotFound
	}

	if order.Status != model.OrderStatusPendingPayment {
		return nil, "", ErrOrderNotPendingPayment
	}

	paymentStatus := model.PaymentStatusSuccess
	orderStatus := "PAID"
	var paidAt any = time.Now()
	now := time.Now()

	if result == model.PaymentStatusFailed {
		paymentStatus = model.PaymentStatusFailed
		orderStatus = "PAYMENT_FAILED"
		paidAt = nil
	}

	paymentNo := generatePaymentNo()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO payments (payment_no, order_id, amount, status, paid_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, paymentNo, order.ID, order.TotalAmount, paymentStatus, paidAt, now, now)
	if err != nil {
		return nil, "", err
	}

	paymentID, err := res.LastInsertId()
	if err != nil {
		return nil, "", err
	}

	if result == model.PaymentStatusFailed {
		type orderItem struct {
			productID int64
			quantity  int64
		}

		var items []orderItem

		rows, err := tx.QueryContext(ctx, `
			SELECT product_id, quantity
			FROM order_items
			WHERE order_id = ?
		`, order.ID)
		if err != nil {
			return nil, "", fmt.Errorf("query order items: %w", err)
		}

		for rows.Next() {
			var item orderItem
			if err := rows.Scan(&item.productID, &item.quantity); err != nil {
				_ = rows.Close()
				return nil, "", fmt.Errorf("scan order item: %w", err)
			}
			items = append(items, item)
		}

		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return nil, "", fmt.Errorf("iterate order items: %w", err)
		}
		if err := rows.Close(); err != nil {
			return nil, "", fmt.Errorf("close order item rows: %w", err)
		}

		for _, item := range items {
			if _, err := tx.ExecContext(ctx, `
				UPDATE inventory
				SET stock = stock + ?, updated_at = ?
				WHERE product_id = ?
			`, item.quantity, now, item.productID); err != nil {
				return nil, "", fmt.Errorf("restore inventory for product %d: %w", item.productID, err)
			}
		}
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE orders
		SET status = ?, updated_at = ?
		WHERE id = ?
	`, orderStatus, now, order.ID); err != nil {
		return nil, "", err
	}

	if err := tx.Commit(); err != nil {
		return nil, "", err
	}

	return &model.Payment{
		ID:        paymentID,
		PaymentNo: paymentNo,
		OrderID:   order.ID,
		Amount:    order.TotalAmount,
		Status:    paymentStatus,
		CreatedAt: now,
		UpdatedAt: now,
	}, orderStatus, nil
}
