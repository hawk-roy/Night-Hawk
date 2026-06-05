package repository

import (
	"context"
	"database/sql"

	"github.com/hawk-roy/Night-Hawk/internal/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			p.id,
			p.name,
			p.description,
			p.price,
			COALESCE(i.stock, 0) AS stock
		FROM products p
		LEFT JOIN inventory i ON p.id = i.product_id
		WHERE p.status = 'ON_SALE'
		ORDER BY p.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
		); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
