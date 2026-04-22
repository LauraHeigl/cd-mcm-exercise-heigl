package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{
		Name:  "Laptop",
		Price: 999.99,
	})

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}
	if got.Name != created.Name {
		t.Errorf("expected Name %q, got %q", created.Name, got.Name)
	}
	if got.Price != created.Price {
		t.Errorf("expected Price %v, got %v", created.Price, got.Price)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestUpdateProduct(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{
		Name:  "Laptop",
		Price: 999.99,
	})

	updatedProduct := model.Product{
		Name:  "Gaming Laptop",
		Price: 1499.99,
	}

	updated, err := s.Update(created.ID, updatedProduct)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, updated.ID)
	}
	if updated.Name != "Gaming Laptop" {
		t.Errorf("expected Name %q, got %q", "Gaming Laptop", updated.Name)
	}
	if updated.Price != 1499.99 {
		t.Errorf("expected Price %v, got %v", 1499.99, updated.Price)
	}

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error when getting updated product, got %v", err)
	}

	if got.Name != "Gaming Laptop" {
		t.Errorf("expected Name %q, got %q", "Gaming Laptop", got.Name)
	}
	if got.Price != 1499.99 {
		t.Errorf("expected Price %v, got %v", 1499.99, got.Price)
	}
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{
		Name:  "Mouse",
		Price: 29.99,
	})

	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()

	tests := []struct {
		name string
		id   int
	}{
		{name: "zero id", id: 0},
		{name: "negative id", id: -1},
		{name: "unknown high id", id: 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetByID(tt.id)
			if err != ErrNotFound {
				t.Errorf("expected ErrNotFound for id %d, got %v", tt.id, err)
			}
		})
	}
}