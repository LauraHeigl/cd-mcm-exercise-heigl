package handler

import (
        "errors"
        "net/http"
        "net/http/httptest"
        "strings"
        "testing"

        "github.com/DATA-DOG/go-sqlmock"
        "github.com/gorilla/mux"
        "github.com/mrckurz/CI-CD-MCM/internal/store"
)

func setupPostgresRouter(t *testing.T) (*mux.Router, sqlmock.Sqlmock, func()) {
        t.Helper()

        db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
        if err != nil {
                t.Fatalf("failed to create sql mock: %v", err)
        }

        pgStore := &store.PostgresStore{DB: db}
        h := NewPostgresHandler(pgStore)
        r := mux.NewRouter()
        h.RegisterRoutes(r)

        cleanup := func() {
                _ = db.Close()
        }

        return r, mock, cleanup
}

func TestPostgresHealthOK(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        mock.ExpectPing()

        req := httptest.NewRequest(http.MethodGet, "/health", nil)
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusOK {
                t.Errorf("expected 200, got %d", rr.Code)
        }
}

func TestPostgresHealthUnavailable(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        mock.ExpectPing().WillReturnError(errors.New("db down"))

        req := httptest.NewRequest(http.MethodGet, "/health", nil)
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusServiceUnavailable {
                t.Errorf("expected 503, got %d", rr.Code)
        }
}

func TestPostgresHandlerGetProducts(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        rows := sqlmock.NewRows([]string{"id", "name", "price"}).
                AddRow(1, "Widget", 9.99)

        mock.ExpectQuery("SELECT id, name, price FROM products ORDER BY id").
                WillReturnRows(rows)

        req := httptest.NewRequest(http.MethodGet, "/products", nil)
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusOK {
                t.Errorf("expected 200, got %d", rr.Code)
        }
}

func TestPostgresHandlerGetProductsError(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        mock.ExpectQuery("SELECT id, name, price FROM products ORDER BY id").
                WillReturnError(errors.New("query failed"))

        req := httptest.NewRequest(http.MethodGet, "/products", nil)
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusInternalServerError {
                t.Errorf("expected 500, got %d", rr.Code)
        }
}

func TestPostgresHandlerGetProduct(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        rows := sqlmock.NewRows([]string{"id", "name", "price"}).
                AddRow(1, "Widget", 9.99)

        mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
                WithArgs(1).
                WillReturnRows(rows)

        req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusOK {
                t.Errorf("expected 200, got %d", rr.Code)
        }
}

func TestPostgresHandlerGetProductNotFound(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        rows := sqlmock.NewRows([]string{"id", "name", "price"})

        mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
                WithArgs(999).
                WillReturnRows(rows)

        req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusNotFound {
                t.Errorf("expected 404, got %d", rr.Code)
        }
}

func TestPostgresHandlerCreateProduct(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

        mock.ExpectQuery("INSERT INTO products").
                WithArgs("Widget", 9.99).
                WillReturnRows(rows)

        req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusCreated {
                t.Errorf("expected 201, got %d", rr.Code)
        }
}

func TestPostgresHandlerCreateProductInvalidJSON(t *testing.T) {
        r, _, cleanup := setupPostgresRouter(t)
        defer cleanup()

        req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{invalid json`))
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusBadRequest {
                t.Errorf("expected 400, got %d", rr.Code)
        }
}

func TestPostgresHandlerCreateProductInvalidProduct(t *testing.T) {
        r, _, cleanup := setupPostgresRouter(t)
        defer cleanup()

        req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name":"","price":-1}`))
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusBadRequest {
                t.Errorf("expected 400, got %d", rr.Code)
        }
}

func TestPostgresHandlerCreateProductStoreError(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        mock.ExpectQuery("INSERT INTO products").
                WithArgs("Widget", 9.99).
                WillReturnError(errors.New("insert failed"))

        req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusInternalServerError {
                t.Errorf("expected 500, got %d", rr.Code)
        }
}

func TestPostgresHandlerUpdateProduct(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        mock.ExpectExec("UPDATE products SET name =").
                WithArgs("Updated", 19.99, 1).
                WillReturnResult(sqlmock.NewResult(0, 1))

        req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(`{"name":"Updated","price":19.99}`))
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusOK {
                t.Errorf("expected 200, got %d", rr.Code)
        }
}

func TestPostgresHandlerUpdateProductInvalidJSON(t *testing.T) {
        r, _, cleanup := setupPostgresRouter(t)
        defer cleanup()

        req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(`{invalid json`))
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusBadRequest {
                t.Errorf("expected 400, got %d", rr.Code)
        }
}

func TestPostgresHandlerUpdateProductNotFound(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        mock.ExpectExec("UPDATE products SET name =").
                WithArgs("Missing", 1.0, 999).
                WillReturnResult(sqlmock.NewResult(0, 0))

        req := httptest.NewRequest(http.MethodPut, "/products/999", strings.NewReader(`{"name":"Missing","price":1}`))
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusNotFound {
                t.Errorf("expected 404, got %d", rr.Code)
        }
}

func TestPostgresHandlerDeleteProduct(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        mock.ExpectExec("DELETE FROM products WHERE id =").
                WithArgs(1).
                WillReturnResult(sqlmock.NewResult(0, 1))

        req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusOK {
                t.Errorf("expected 200, got %d", rr.Code)
        }
}

func TestPostgresHandlerDeleteProductNotFound(t *testing.T) {
        r, mock, cleanup := setupPostgresRouter(t)
        defer cleanup()

        mock.ExpectExec("DELETE FROM products WHERE id =").
                WithArgs(999).
                WillReturnResult(sqlmock.NewResult(0, 0))

        req := httptest.NewRequest(http.MethodDelete, "/products/999", nil)
        rr := httptest.NewRecorder()

        r.ServeHTTP(rr, req)

        if rr.Code != http.StatusNotFound {
                t.Errorf("expected 404, got %d", rr.Code)
        }
}