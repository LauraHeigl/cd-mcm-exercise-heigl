package store

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func newMockPostgresStore(t *testing.T) (*PostgresStore, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	cleanup := func() {
		_ = db.Close()
	}

	return &PostgresStore{DB: db}, mock, cleanup
}

func TestPostgresEnsureTable(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS products").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := s.EnsureTable(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPostgresEnsureTableError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS products").
		WillReturnError(errors.New("create table failed"))

	if err := s.EnsureTable(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPostgresGetAll(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Product A", 10.0).
		AddRow(2, "Product B", 20.0)

	mock.ExpectQuery("SELECT id, name, price FROM products").
		WillReturnRows(rows)

	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
}

func TestPostgresGetAllQueryError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, name, price FROM products").
		WillReturnError(errors.New("query failed"))

	_, err := s.GetAll()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPostgresGetByID(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Product A", 10.0)

	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(1).
		WillReturnRows(rows)

	product, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != 1 {
		t.Errorf("expected ID 1, got %d", product.ID)
	}
}

func TestPostgresGetByIDNotFound(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err := s.GetByID(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresGetByIDDatabaseError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err := s.GetByID(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPostgresCreate(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO products").
		WithArgs("Product A", 10.0).
		WillReturnRows(rows)

	product, err := s.Create(model.Product{
		Name:  "Product A",
		Price: 10.0,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != 1 {
		t.Errorf("expected ID 1, got %d", product.ID)
	}
}

func TestPostgresCreateError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectQuery("INSERT INTO products").
		WithArgs("Product A", 10.0).
		WillReturnError(errors.New("insert failed"))

	_, err := s.Create(model.Product{
		Name:  "Product A",
		Price: 10.0,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPostgresUpdate(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("UPDATE products SET name").
		WithArgs("Updated", 15.0, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	product, err := s.Update(1, model.Product{
		Name:  "Updated",
		Price: 15.0,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != 1 {
		t.Errorf("expected ID 1, got %d", product.ID)
	}
}

func TestPostgresUpdateNotFound(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("UPDATE products SET name").
		WithArgs("Missing", 15.0, 999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := s.Update(999, model.Product{
		Name:  "Missing",
		Price: 15.0,
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresUpdateError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("UPDATE products SET name").
		WithArgs("Updated", 15.0, 1).
		WillReturnError(errors.New("update failed"))

	_, err := s.Update(1, model.Product{
		Name:  "Updated",
		Price: 15.0,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPostgresDelete(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM products WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := s.Delete(1); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM products WHERE id").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := s.Delete(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresDeleteError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM products WHERE id").
		WithArgs(1).
		WillReturnError(errors.New("delete failed"))

	err := s.Delete(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
