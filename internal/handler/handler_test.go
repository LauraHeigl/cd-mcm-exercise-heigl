package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
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

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductsEmpty(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestCreateAndGetProduct(t *testing.T) {
	r, _ := setupRouter()

	// Create
	body := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	// Get
	req = httptest.NewRequest("GET", "/products/1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductNotFound(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUpdateProduct(t *testing.T) {
router, _ := setupRouter()

	createBody := strings.NewReader(`{"name":"Laptop","price":999.99}`)
	createReq := httptest.NewRequest(http.MethodPost, "/products", createBody)
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createRec.Code)
	}

	updateBody := strings.NewReader(`{"name":"Gaming Laptop","price":1299.99}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/products/1", updateBody)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()

	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", updateRec.Code)
	}

	if !strings.Contains(updateRec.Body.String(), "Gaming Laptop") {
		t.Errorf("expected updated product name in response, got %s", updateRec.Body.String())
	}
}

func TestDeleteProduct(t *testing.T) {
router, _ := setupRouter()

	createBody := strings.NewReader(`{"name":"Laptop","price":999.99}`)
	createReq := httptest.NewRequest(http.MethodPost, "/products", createBody)
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createRec.Code)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	deleteRec := httptest.NewRecorder()

	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", deleteRec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	getRec := httptest.NewRecorder()

	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 after delete, got %d", getRec.Code)
	}
}

func TestCreateInvalidProduct(t *testing.T) {
router, _ := setupRouter()

	body := strings.NewReader(`{"name":"","price":999.99}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

// TODO: Add tests for UpdateProduct, DeleteProduct, and invalid payloads
