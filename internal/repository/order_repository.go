package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hawk-roy/Night-Hawk/internal/model"
)

var ErrProductNotFound = errors.New("product not found")
var ErrInsufficientStock = errors.New("insufficient stock")

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func generateOrderNo() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}

func (r *OrderRepository) CreateOrder(ctx context.Context, userID int64, username string, productID int64, quantity int64) (*model.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var product model.Product
	var status string

	row := tx.QueryRowContext(ctx, `
		SELECT
			p.id,
			p.name,
			p.price,
			p.status,
			i.stock
		FROM products p
		JOIN inventory i ON p.id = i.product_id
		WHERE p.id = ?
		FOR UPDATE
	`, productID)

	if err := row.Scan(&product.ID, &product.Name, &product.Price, &status, &product.Stock); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if status != "ON_SALE" {
		return nil, ErrProductNotFound
	}

	if product.Stock < quantity {
		return nil, ErrInsufficientStock
	}

	orderNo := generateOrderNo()
	totalAmount := product.Price * quantity
	now := time.Now()

	result, err := tx.ExecContext(ctx, `
		UPDATE inventory
		SET stock = stock - ?, updated_at = ?
		WHERE product_id = ? AND stock >= ?
	`, quantity, now, productID, quantity)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrInsufficientStock
	}

	result, err = tx.ExecContext(ctx, `
		INSERT INTO orders (order_no, user_id, total_amount, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, orderNo, userID, totalAmount, model.OrderStatusPendingPayment, now, now)
	if err != nil {
		return nil, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO order_items (order_id, product_id, product_name, price, quantity, subtotal, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, orderID, product.ID, product.Name, product.Price, quantity, totalAmount, now, now)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Order{
		ID:          orderID,
		OrderNo:     orderNo,
		UserID:      userID,
		Username:    username,
		ProductID:   product.ID,
		ProductName: product.Name,
		UnitPrice:   product.Price,
		Quantity:    quantity,
		TotalAmount: totalAmount,
		Status:      model.OrderStatusPendingPayment,
		CreatedAt:   now,
	}, nil
}
