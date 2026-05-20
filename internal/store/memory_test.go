package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	p := model.Product{Name: "Test", Price: 1.99}
	created := s.Create(p)
	if created.ID == 0 {
		t.Error("expected created product to have an ID")
	}
	if created.Name != p.Name || created.Price != p.Price {
		t.Errorf("expected created product to match input, got %+v", created)
	}

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got.ID != created.ID || got.Name != created.Name || got.Price != created.Price {
		t.Errorf("expected to get the same product back, got %+v", got)
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

func TestUpdate(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(model.Product{Name: "Test", Price: 1.99})
	p.Name = "Updated"
	updated, err := s.Update(p.ID, p)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("expected name to be 'Updated', got '%s'", updated.Name)
	}
}

func TestDeleteExisting(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(model.Product{Name: "Test", Price: 1.99})
	err := s.Delete(p.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	_, err = s.GetByID(p.ID)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound after deleting product")
	}
}

func TestGetByID(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(model.Product{Name: "Test", Price: 1.99})
	got, err := s.GetByID(p.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got.ID != p.ID || got.Name != p.Name || got.Price != p.Price {
		t.Errorf("expected to get the same product back, got %+v", got)
	}
}

func TestGetByIDInvalid(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when getting non-existent product")
	}
}
