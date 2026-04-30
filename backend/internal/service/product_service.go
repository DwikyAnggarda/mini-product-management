package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"product-management/backend/internal/model"
	"product-management/backend/internal/repository"
)

type ProductService interface {
	List(ctx context.Context, query, status, pageStr, limitStr string) ([]model.Product, map[string]interface{}, error)
	Create(ctx context.Context, payload model.ProductPayload) (model.Product, map[string]string, error)
	Update(ctx context.Context, id int64, payload model.ProductPayload) (model.Product, map[string]string, error)
	Delete(ctx context.Context, id int64) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) List(ctx context.Context, query, status, pageStr, limitStr string) ([]model.Product, map[string]interface{}, error) {
	page := 1
	limit := 10
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p <= 0 {
			return nil, nil, fmt.Errorf("page harus bilangan bulat positif")
		}
		page = p
	}
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 || l > 100 {
			return nil, nil, fmt.Errorf("limit harus 1-100")
		}
		limit = l
	}

	status = strings.TrimSpace(status)
	if status != "" && status != "active" && status != "inactive" {
		return nil, nil, fmt.Errorf("status filter harus active atau inactive")
	}

	params := model.ProductListParams{
		Query:  strings.TrimSpace(query),
		Status: status,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}
	meta := map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	}
	return items, meta, nil
}

func (s *productService) Create(ctx context.Context, payload model.ProductPayload) (model.Product, map[string]string, error) {
	payload.SKU = strings.TrimSpace(payload.SKU)
	payload.Name = strings.TrimSpace(payload.Name)
	payload.Status = strings.TrimSpace(payload.Status)

	if errs := payload.Validate(); len(errs) > 0 {
		return model.Product{}, errs, fmt.Errorf("validation error")
	}

	product, err := s.repo.Create(ctx, payload)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return model.Product{}, map[string]string{"sku": "sku sudah terdaftar"}, fmt.Errorf("validation error")
		}
		return model.Product{}, nil, err
	}
	return product, nil, nil
}

func (s *productService) Update(ctx context.Context, id int64, payload model.ProductPayload) (model.Product, map[string]string, error) {
	payload.SKU = strings.TrimSpace(payload.SKU)
	payload.Name = strings.TrimSpace(payload.Name)
	payload.Status = strings.TrimSpace(payload.Status)

	if errs := payload.Validate(); len(errs) > 0 {
		return model.Product{}, errs, fmt.Errorf("validation error")
	}

	product, err := s.repo.Update(ctx, id, payload)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return model.Product{}, map[string]string{"sku": "sku sudah terdaftar"}, fmt.Errorf("validation error")
		}
		return model.Product{}, nil, err
	}
	return product, nil, nil
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
