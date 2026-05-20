package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func setupRouter() (*mux.Router, *Handler) {
	s := store.NewMemoryStore()
	h := NewHandler(s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r, h
}

func TestHealthEndpoint(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductsEmpty(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var products []model.Product
	if err := json.NewDecoder(rr.Body).Decode(&products); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(products) != 0 {
		t.Errorf("expected empty product list, got %d products", len(products))
	}
}

func TestCreateAndGetProduct(t *testing.T) {
	r, _ := setupRouter()

	body := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestCreateProductInvalidJSON(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{invalid json`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateProductInvalidProduct(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name":"","price":-1}`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestGetProductNotFound(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUpdateProduct(t *testing.T) {
	r, _ := setupRouter()

	createBody := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(createBody))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	updateBody := `{"name":"Updated Widget","price":19.99}`
	req = httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(updateBody))
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var updated model.Product
	if err := json.NewDecoder(rr.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.ID != 1 {
		t.Errorf("expected ID 1, got %d", updated.ID)
	}

	if updated.Name != "Updated Widget" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}

	if updated.Price != 19.99 {
		t.Errorf("expected updated price 19.99, got %f", updated.Price)
	}
}

func TestUpdateProductInvalidJSON(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(`{invalid json`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateProductNotFound(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodPut, "/products/999", strings.NewReader(`{"name":"Missing","price":1}`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestDeleteProduct(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	req = httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected deleted product to return 404, got %d", rr.Code)
	}
}

func TestDeleteProductNotFound(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/products/999", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
