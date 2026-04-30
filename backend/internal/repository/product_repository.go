package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"product-management/backend/internal/model"
)

type ProductRepository interface {
	List(ctx context.Context, params model.ProductListParams) ([]model.Product, int, error)
	Create(ctx context.Context, payload model.ProductPayload) (model.Product, error)
	GetByID(ctx context.Context, id int64) (model.Product, error)
	Update(ctx context.Context, id int64, payload model.ProductPayload) (model.Product, error)
	Delete(ctx context.Context, id int64) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) List(ctx context.Context, params model.ProductListParams) ([]model.Product, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if params.Query != "" {
		args = append(args, "%"+params.Query+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR sku ILIKE $%d)", idx, idx))
	}
	if params.Status != "" {
		args = append(args, params.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, params.Limit, params.Offset)
	limitArg := len(args) - 1
	offsetArg := len(args)
	listQuery := fmt.Sprintf(`
		SELECT id, sku, name, description, price, status, created_at, updated_at
		FROM products
		WHERE %s
		ORDER BY updated_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, limitArg, offsetArg)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Create(ctx context.Context, payload model.ProductPayload) (model.Product, error) {
	query := `
		INSERT INTO products (sku, name, description, price, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, sku, name, description, price, status, created_at, updated_at`
	var p model.Product
	err := r.db.QueryRowContext(ctx, query, payload.SKU, payload.Name, payload.Description, payload.Price, payload.Status).
		Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return model.Product{}, err
	}
	return p, nil
}

func (r *productRepository) GetByID(ctx context.Context, id int64) (model.Product, error) {
	query := `SELECT id, sku, name, description, price, status, created_at, updated_at FROM products WHERE id = $1`
	var p model.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, fmt.Errorf("product not found")
		}
		return model.Product{}, err
	}
	return p, nil
}

func (r *productRepository) Update(ctx context.Context, id int64, payload model.ProductPayload) (model.Product, error) {
	query := `
		UPDATE products
		SET sku = $2, name = $3, description = $4, price = $5, status = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, sku, name, description, price, status, created_at, updated_at`
	var p model.Product
	err := r.db.QueryRowContext(ctx, query, id, payload.SKU, payload.Name, payload.Description, payload.Price, payload.Status).
		Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, fmt.Errorf("product not found")
		}
		return model.Product{}, err
	}
	return p, nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
