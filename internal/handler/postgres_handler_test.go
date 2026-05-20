package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func setupPostgresRouter(t *testing.T) (*mux.Router, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	s := &store.PostgresStore{DB: db}
	h := NewPostgresHandler(s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r, mock, db
}

func TestPostgresHealthOK(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	mock.ExpectPing()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProducts(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresCreateProduct(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO products").
		WithArgs("Widget", 9.99).WillReturnRows(rows)

	body := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestPostgresGetProductNotFound(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(999).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}))

	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresGetProductOK(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(1).WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresCreateInvalidJSON(t *testing.T) {
	r, _, db := setupPostgresRouter(t)
	defer db.Close()

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{bad}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresCreateInvalidProduct(t *testing.T) {
	r, _, db := setupPostgresRouter(t)
	defer db.Close()

	body := `{"price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresUpdateProduct(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	mock.ExpectExec("UPDATE products").
		WithArgs("Updated", 19.99, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"name":"Updated","price":19.99}`
	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresDeleteProduct(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM products").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresUpdateInvalidJSON(t *testing.T) {
	r, _, db := setupPostgresRouter(t)
	defer db.Close()

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{bad}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresUpdateNotFound(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	mock.ExpectExec("UPDATE products").
		WithArgs("X", 1.0, 999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	body := `{"name":"X","price":1.0}`
	req := httptest.NewRequest("PUT", "/products/999", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM products").
		WithArgs(999).WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest("DELETE", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresHealthOKPing(t *testing.T) {
	r, mock, db := setupPostgresRouter(t)
	defer db.Close()

	mock.ExpectPing()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
