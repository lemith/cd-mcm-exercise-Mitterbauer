//go:build !unit

package store

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func newMockStore(t *testing.T) (*PostgresStore, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return &PostgresStore{DB: db}, mock
}

func TestPostgresGetAll(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Widget", 9.99).
		AddRow(2, "Gadget", 19.99)
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(rows)

	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestPostgresGetByID(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(1).WillReturnRows(rows)

	p, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Widget" {
		t.Errorf("expected Widget, got %s", p.Name)
	}
}

func TestPostgresGetByIDNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price"})
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(999).WillReturnRows(rows)

	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound")
	}
}

func TestPostgresCreate(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO products").
		WithArgs("Widget", 9.99).WillReturnRows(rows)

	p, err := s.Create(model.Product{Name: "Widget", Price: 9.99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}
}

func TestPostgresUpdate(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	mock.ExpectExec("UPDATE products").
		WithArgs("Updated", 19.99, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	p, err := s.Update(1, model.Product{Name: "Updated", Price: 19.99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Updated" {
		t.Errorf("expected Updated, got %s", p.Name)
	}
}

func TestPostgresUpdateNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	mock.ExpectExec("UPDATE products").
		WithArgs("X", 1.0, 999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := s.Update(999, model.Product{Name: "X", Price: 1.0})
	if err != ErrNotFound {
		t.Error("expected ErrNotFound")
	}
}

func TestPostgresDelete(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	mock.ExpectExec("DELETE FROM products").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.Delete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	mock.ExpectExec("DELETE FROM products").
		WithArgs(999).WillReturnResult(sqlmock.NewResult(0, 0))

	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound")
	}
}

func TestPostgresEnsureTable(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS products").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := s.EnsureTable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewPostgresStoreFailure(t *testing.T) {
	_, err := NewPostgresStore("invalid", "0", "user", "pass", "db")
	if err == nil {
		t.Error("expected error for invalid connection")
	}
}

func TestPostgresGetAllError(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	mock.ExpectQuery("SELECT id, name, price FROM products").
		WillReturnError(fmt.Errorf("db error"))

	_, err := s.GetAll()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPostgresUpdateDBError(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	mock.ExpectExec("UPDATE products").
		WithArgs("X", 1.0, 1).
		WillReturnError(fmt.Errorf("db error"))

	_, err := s.Update(1, model.Product{Name: "X", Price: 1.0})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPostgresDeleteDBError(t *testing.T) {
	s, mock := newMockStore(t)
	defer s.DB.Close()

	mock.ExpectExec("DELETE FROM products").
		WithArgs(1).
		WillReturnError(fmt.Errorf("db error"))

	err := s.Delete(1)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
