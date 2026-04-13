package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{})
	p, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, p.ID)
	}
	// TODO: Add test -- create a product and verify GetByID returns it
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

// TODO: Add tests for Update, Delete of existing product, and GetByID with invalid ID
func TestUpdateProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{})
	created.Name = "Updated Name"
	updated, err := s.Update(created.ID, created)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated, err = s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", updated.Name)
	}
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{})
	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound after deleting product")
	}
}

// func TestGetByIDNotFound(t *testing.T) {
//  s := NewMemoryStore()
//  _, err := s.GetByID(999)
//  if err != ErrNotFound {
//      t.Error("expected ErrNotFound when getting non-existent product")
//  }
// }

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Test Product"})

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{"existing product", created.ID, nil},
		{"non-existent", 999, ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetByID(tt.id)
			if err != tt.wantErr {
				t.Errorf("GetByID(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}
