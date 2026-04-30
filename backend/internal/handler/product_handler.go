package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"product-management/backend/internal/model"
	"product-management/backend/internal/response"
	"product-management/backend/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	status := r.URL.Query().Get("status")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	products, meta, err := h.productService.List(r.Context(), query, status, page, limit)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, products, meta)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload model.ProductPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body", nil)
		return
	}

	product, details, err := h.productService.Create(r.Context(), payload)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "input validation failed", details)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create product", nil)
		return
	}

	response.JSON(w, http.StatusCreated, product, nil)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	var payload model.ProductPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body", nil)
		return
	}

	product, details, svcErr := h.productService.Update(r.Context(), id, payload)
	if svcErr != nil {
		if strings.Contains(svcErr.Error(), "validation") {
			response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "input validation failed", details)
			return
		}
		if strings.Contains(strings.ToLower(svcErr.Error()), "not found") {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "product not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update product", nil)
		return
	}

	response.JSON(w, http.StatusOK, product, nil)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	if err := h.productService.Delete(r.Context(), id); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "product not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete product", nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "product deleted"}, nil)
}

func parseID(idRaw string) (int64, error) {
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}
