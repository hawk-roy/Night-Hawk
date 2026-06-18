package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/hawk-roy/Night-Hawk/internal/model"
)

func TestPaymentRepository_MockPayOrderSuccess(t *testing.T) {
	skipIntegrationTest(t)

	db := openIntegrationDB(t)

	userID, username := createTestUser(t, db, "it_payment_user")
	productID, _ := createTestProductWithInventory(t, db, "it_payment_product", 10, 19900)

	orderRepo := NewOrderRepository(db)
	order, err := orderRepo.CreateOrder(context.Background(), userID, username, productID, 2)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	var paymentID int64
	t.Cleanup(func() {
		cleanupTestData(t, db, userID, productID, order.ID, paymentID)
	})

	if order.Status != model.OrderStatusPendingPayment {
		t.Fatalf("unexpected order status before payment: got %q want %q", order.Status, model.OrderStatusPendingPayment)
	}

	paymentRepo := NewPaymentRepository(db)
	payment, orderStatus, err := paymentRepo.MockPayOrder(context.Background(), userID, order.ID, model.PaymentStatusSuccess)
	if err != nil {
		t.Fatalf("MockPayOrder success failed: %v", err)
	}
	if payment == nil {
		t.Fatal("expected payment to be returned")
	}

	paymentID = payment.ID

	if payment.Status != model.PaymentStatusSuccess {
		t.Fatalf("unexpected payment status: got %q want %q", payment.Status, model.PaymentStatusSuccess)
	}
	if orderStatus != "PAID" {
		t.Fatalf("unexpected order status returned: got %q want %q", orderStatus, "PAID")
	}

	var paymentCount int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT COUNT(*)
		FROM payments
		WHERE order_id = ?
	`, order.ID).Scan(&paymentCount); err != nil {
		t.Fatalf("failed to count payments: %v", err)
	}
	if paymentCount != 1 {
		t.Fatalf("unexpected payment count: got %d want %d", paymentCount, 1)
	}

	var dbPaymentStatus string
	var dbPaymentAmount int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT status, amount
		FROM payments
		WHERE order_id = ?
		ORDER BY id DESC
		LIMIT 1
	`, order.ID).Scan(&dbPaymentStatus, &dbPaymentAmount); err != nil {
		t.Fatalf("failed to query payment: %v", err)
	}
	if dbPaymentStatus != model.PaymentStatusSuccess {
		t.Fatalf("unexpected db payment status: got %q want %q", dbPaymentStatus, model.PaymentStatusSuccess)
	}
	if dbPaymentAmount != 39800 {
		t.Fatalf("unexpected db payment amount: got %d want %d", dbPaymentAmount, 39800)
	}

	var dbOrderStatus string
	if err := db.QueryRowContext(context.Background(), `
		SELECT status
		FROM orders
		WHERE id = ?
	`, order.ID).Scan(&dbOrderStatus); err != nil {
		t.Fatalf("failed to query order status: %v", err)
	}
	if dbOrderStatus != "PAID" {
		t.Fatalf("unexpected db order status: got %q want %q", dbOrderStatus, "PAID")
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

	_, _, err = paymentRepo.MockPayOrder(context.Background(), userID, order.ID, model.PaymentStatusSuccess)
	if !errors.Is(err, ErrOrderNotPendingPayment) {
		t.Fatalf("unexpected repeat payment error: got %v want %v", err, ErrOrderNotPendingPayment)
	}
}

func TestPaymentRepository_MockPayOrderFailedRestoresStock(t *testing.T) {
	skipIntegrationTest(t)

	db := openIntegrationDB(t)

	userID, username := createTestUser(t, db, "it_payment_user")
	productID, _ := createTestProductWithInventory(t, db, "it_payment_product", 10, 19900)

	orderRepo := NewOrderRepository(db)
	order, err := orderRepo.CreateOrder(context.Background(), userID, username, productID, 3)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	var paymentID int64
	t.Cleanup(func() {
		cleanupTestData(t, db, userID, productID, order.ID, paymentID)
	})

	var stockBeforePayment int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT stock
		FROM inventory
		WHERE product_id = ?
	`, productID).Scan(&stockBeforePayment); err != nil {
		t.Fatalf("failed to query inventory before payment: %v", err)
	}
	if stockBeforePayment != 7 {
		t.Fatalf("unexpected inventory stock before payment: got %d want %d", stockBeforePayment, 7)
	}

	paymentRepo := NewPaymentRepository(db)
	payment, orderStatus, err := paymentRepo.MockPayOrder(context.Background(), userID, order.ID, model.PaymentStatusFailed)
	if err != nil {
		t.Fatalf("MockPayOrder failed result failed: %v", err)
	}
	if payment == nil {
		t.Fatal("expected payment to be returned")
	}

	paymentID = payment.ID

	if payment.Status != model.PaymentStatusFailed {
		t.Fatalf("unexpected payment status: got %q want %q", payment.Status, model.PaymentStatusFailed)
	}
	if orderStatus != "PAYMENT_FAILED" {
		t.Fatalf("unexpected order status returned: got %q want %q", orderStatus, "PAYMENT_FAILED")
	}

	var paymentCount int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT COUNT(*)
		FROM payments
		WHERE order_id = ?
	`, order.ID).Scan(&paymentCount); err != nil {
		t.Fatalf("failed to count payments: %v", err)
	}
	if paymentCount != 1 {
		t.Fatalf("unexpected payment count: got %d want %d", paymentCount, 1)
	}

	var dbPaymentStatus string
	var dbPaymentAmount int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT status, amount
		FROM payments
		WHERE order_id = ?
		ORDER BY id DESC
		LIMIT 1
	`, order.ID).Scan(&dbPaymentStatus, &dbPaymentAmount); err != nil {
		t.Fatalf("failed to query payment: %v", err)
	}
	if dbPaymentStatus != model.PaymentStatusFailed {
		t.Fatalf("unexpected db payment status: got %q want %q", dbPaymentStatus, model.PaymentStatusFailed)
	}
	if dbPaymentAmount != 59700 {
		t.Fatalf("unexpected db payment amount: got %d want %d", dbPaymentAmount, 59700)
	}

	var dbOrderStatus string
	if err := db.QueryRowContext(context.Background(), `
		SELECT status
		FROM orders
		WHERE id = ?
	`, order.ID).Scan(&dbOrderStatus); err != nil {
		t.Fatalf("failed to query order status: %v", err)
	}
	if dbOrderStatus != "PAYMENT_FAILED" {
		t.Fatalf("unexpected db order status: got %q want %q", dbOrderStatus, "PAYMENT_FAILED")
	}

	var stockAfterPayment int64
	if err := db.QueryRowContext(context.Background(), `
		SELECT stock
		FROM inventory
		WHERE product_id = ?
	`, productID).Scan(&stockAfterPayment); err != nil {
		t.Fatalf("failed to query inventory after payment: %v", err)
	}
	if stockAfterPayment != 10 {
		t.Fatalf("unexpected inventory stock after failed payment: got %d want %d", stockAfterPayment, 10)
	}
}
