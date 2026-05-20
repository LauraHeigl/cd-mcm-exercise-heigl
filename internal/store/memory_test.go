package store

import (
	"errors"
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{
		Name:  "Widget",
		Price: 9.99,
	})

	if created.ID != 1 {
		t.Fatalf("expected ID 1, got %d", created.ID)
	}

	found, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected product, got error: %v", err)
	}

	if found.Name != "Widget" {
		t.Errorf("expected name Widget, got %s", found.Name)
	}

	if found.Price != 9.99 {
		t.Errorf("expected price 9.99, got %f", found.Price)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()

	products := s.GetAll()

	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestGetAllWithProducts(t *testing.T) {
	s := NewMemoryStore()

	s.Create(model.Product{Name: "A", Price: 1})
	s.Create(model.Product{Name: "B", Price: 2})

	products := s.GetAll()

	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.GetByID(999)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateExistingProduct(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{Name: "Old", Price: 5})

	updated, err := s.Update(created.ID, model.Product{
		Name:  "New",
		Price: 10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, updated.ID)
	}

	if updated.Name != "New" {
		t.Errorf("expected name New, got %s", updated.Name)
	}

	found, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected product after update, got %v", err)
	}

	if found.Name != "New" {
		t.Errorf("expected stored product to be updated, got %s", found.Name)
	}
}

func TestUpdateNonExistentProduct(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.Update(999, model.Product{Name: "Missing", Price: 1})

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteExistingProduct(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{Name: "Delete me", Price: 3})

	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.GetByID(created.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected product to be deleted, got %v", err)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()

	err := s.Delete(999)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound when deleting non-existent product, got %v", err)
	}
}
