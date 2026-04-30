package model

import "time"

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type Product struct {
	ID          int64     `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductPayload struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
}

func (p ProductPayload) Validate() map[string]string {
	errs := map[string]string{}
	if len(p.SKU) < 3 || len(p.SKU) > 40 {
		errs["sku"] = "sku harus 3-40 karakter"
	}
	if len(p.Name) < 2 || len(p.Name) > 120 {
		errs["name"] = "name harus 2-120 karakter"
	}
	if p.Price < 0 {
		errs["price"] = "price tidak boleh negatif"
	}
	if p.Status != "active" && p.Status != "inactive" {
		errs["status"] = "status harus active atau inactive"
	}
	return errs
}

type ProductListParams struct {
	Query  string
	Status string
	Limit  int
	Offset int
}
