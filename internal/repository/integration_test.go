package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func skipIntegrationTest(t *testing.T) {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("skip integration test: set RUN_INTEGRATION_TESTS=1 to run")
	}
}

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	host := getenvDefault("MYSQL_HOST", "127.0.0.1")
	port := getenvDefault("MYSQL_PORT", "3306")
	database := getenvDefault("MYSQL_DATABASE", "go_order_service")
	user := getenvDefault("MYSQL_USER", "order_user")
	password := getenvDefault("MYSQL_PASSWORD", "order_pass")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		user,
		password,
		host,
		port,
		database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open mysql: %v", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		t.Fatalf("failed to ping mysql: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func createTestUser(t *testing.T, db *sql.DB, prefix string) (int64, string) {
	t.Helper()

	if prefix == "" {
		prefix = "it_user"
	}

	username := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	passwordHash := "hash_" + username

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.ExecContext(ctx, `
		INSERT INTO users (username, password_hash, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
	`, username, passwordHash)
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get test user id: %v", err)
	}

	return userID, username
}

func createTestProductWithInventory(t *testing.T, db *sql.DB, prefix string, stock int64, price int64) (int64, string) {
	t.Helper()

	if prefix == "" {
		prefix = "it_product"
	}

	productName := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	description := "integration test product"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("failed to begin test tx: %v", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO products (name, description, price, status, created_at, updated_at)
		VALUES (?, ?, ?, 'ON_SALE', NOW(), NOW())
	`, productName, description, price)
	if err != nil {
		t.Fatalf("failed to insert test product: %v", err)
	}

	productID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get test product id: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO inventory (product_id, stock, locked_stock, created_at, updated_at)
		VALUES (?, ?, 0, NOW(), NOW())
	`, productID, stock)
	if err != nil {
		t.Fatalf("failed to insert test inventory: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("failed to commit test product tx: %v", err)
	}

	return productID, productName
}

func cleanupTestData(t *testing.T, db *sql.DB, userID, productID, orderID, paymentID int64) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if paymentID > 0 {
		if _, err := db.ExecContext(ctx, `DELETE FROM payments WHERE id = ?`, paymentID); err != nil {
			t.Fatalf("failed to delete test payment: %v", err)
		}
	}

	if orderID > 0 {
		if _, err := db.ExecContext(ctx, `DELETE FROM order_items WHERE order_id = ?`, orderID); err != nil {
			t.Fatalf("failed to delete test order items: %v", err)
		}
		if _, err := db.ExecContext(ctx, `DELETE FROM orders WHERE id = ?`, orderID); err != nil {
			t.Fatalf("failed to delete test order: %v", err)
		}
	}

	if productID > 0 {
		if _, err := db.ExecContext(ctx, `DELETE FROM inventory WHERE product_id = ?`, productID); err != nil {
			t.Fatalf("failed to delete test inventory: %v", err)
		}
		if _, err := db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, productID); err != nil {
			t.Fatalf("failed to delete test product: %v", err)
		}
	}

	if userID > 0 {
		if _, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, userID); err != nil {
			t.Fatalf("failed to delete test user: %v", err)
		}
	}
}

func getenvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
