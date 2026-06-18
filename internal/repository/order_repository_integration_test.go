package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/hawk-roy/Night-Hawk/internal/model"
)

func TestOrderRepository_CreateOrder_DeductsStock(t *testing.T) {
	skipIntegrationTest(t)

	db := openIntegrationDB(t)

	userID, username := createTestUser(t, db, "it_order_user")
	productID, _ := createTestProductWithInventory(t, db, "it_order_product", 10, 19900)
	defer cleanupTestData(t, db, userID, productID, 0, 0)

	repo := NewOrderRepository(db)

	order, err := repo.CreateOrder(context.Background(), userID, username, productID, 2)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if order == nil {
		t.Fatal("expected order to be returned")
	}

	if order.Status != model.OrderStatusPendingPayment {
		t.Fatalf("unexpected order status: got %q want %q", order.Status, model.OrderStatusPendingPayment)
	}
	if order.ProductID != productID {
		t.Fatalf("unexpected product id: got %d want %d", order.ProductID, productID)
	}
	if order.Quantity != 2 {
		t.Fatalf("unexpected quantity: got %d want %d", order.Quantity, 2)
	}
	if order.TotalAmount != 39800 {
		t.Fatalf("unexpected total amount: got %d want %d", order.TotalAmount, 39800)
	}

	var dbStatus string
	var dbTotalAmount int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT status, total_amount
		FROM orders
		WHERE id = ?
	`, order.ID).Scan(&dbStatus, &dbTotalAmount); err != nil {
		t.Fatalf("failed to query order: %v", err)
	}

	if dbStatus != model.OrderStatusPendingPayment {
		t.Fatalf("unexpected db order status: got %q want %q", dbStatus, model.OrderStatusPendingPayment)
	}
	if dbTotalAmount != 39800 {
		t.Fatalf("unexpected db total amount: got %d want %d", dbTotalAmount, 39800)
	}

	var itemCount int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT COUNT(*)
		FROM order_items
		WHERE order_id = ?
	`, order.ID).Scan(&itemCount); err != nil {
		t.Fatalf("failed to count order items: %v", err)
	}
	if itemCount != 1 {
		t.Fatalf("unexpected order item count: got %d want %d", itemCount, 1)
	}

	var itemProductID, itemQuantity int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT product_id, quantity
		FROM order_items
		WHERE order_id = ?
	`, order.ID).Scan(&itemProductID, &itemQuantity); err != nil {
		t.Fatalf("failed to query order item: %v", err)
	}

	if itemProductID != productID {
		t.Fatalf("unexpected order item product id: got %d want %d", itemProductID, productID)
	}
	if itemQuantity != 2 {
		t.Fatalf("unexpected order item quantity: got %d want %d", itemQuantity, 2)
	}

	var stock int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT stock
		FROM inventory
		WHERE product_id = ?
	`, productID).Scan(&stock); err != nil {
		t.Fatalf("failed to query inventory: %v", err)
	}
	if stock != 8 {
		t.Fatalf("unexpected inventory stock: got %d want %d", stock, 8)
	}
}

func TestOrderRepository_CreateOrder_InsufficientStockRollback(t *testing.T) {
	skipIntegrationTest(t)

	db := openIntegrationDB(t)

	userID, username := createTestUser(t, db, "it_order_user")
	productID, _ := createTestProductWithInventory(t, db, "it_order_product", 1, 19900)
	defer cleanupTestData(t, db, userID, productID, 0, 0)

	repo := NewOrderRepository(db)

	order, err := repo.CreateOrder(context.Background(), userID, username, productID, 2)
	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("unexpected error: got %v want %v", err, ErrInsufficientStock)
	}
	if order != nil {
		t.Fatal("expected no order to be returned on insufficient stock")
	}

	var stock int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT stock
		FROM inventory
		WHERE product_id = ?
	`, productID).Scan(&stock); err != nil {
		t.Fatalf("failed to query inventory: %v", err)
	}
	if stock != 1 {
		t.Fatalf("unexpected inventory stock: got %d want %d", stock, 1)
	}

	var orderCount int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = ?
	`, userID).Scan(&orderCount); err != nil {
		t.Fatalf("failed to count orders: %v", err)
	}
	if orderCount != 0 {
		t.Fatalf("unexpected order count: got %d want %d", orderCount, 0)
	}

	var itemCount int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT COUNT(*)
		FROM order_items
		WHERE product_id = ?
	`, productID).Scan(&itemCount); err != nil {
		t.Fatalf("failed to count order items: %v", err)
	}
	if itemCount != 0 {
		t.Fatalf("unexpected order item count: got %d want %d", itemCount, 0)
	}
}
