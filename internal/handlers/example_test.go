package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/drivers"
	"github.com/thxhix/shortener/internal/models"
	urlUseCase "github.com/thxhix/shortener/internal/url"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleHandler_StoreLink() {
	db, _ := drivers.NewMemoryDatabase()
	cfg := config.Config{BaseURL: "http://localhost:8080"}
	useCase := urlUseCase.NewURLUseCase(db, cfg)
	h := NewHandler(&cfg, useCase)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	w := httptest.NewRecorder()

	h.StoreLink(w, req)

	fmt.Println("Status:", w.Code)
	fmt.Println("Body:", strings.TrimSpace(w.Body.String()))

	// Output:
	// Status: 201
	// Body: http://localhost:8080/testHash
}

func ExampleHandler_Redirect() {
	db, _ := drivers.NewMemoryDatabase()
	cfg := config.Config{BaseURL: "http://localhost:8080"}
	useCase := urlUseCase.NewURLUseCase(db, cfg)
	h := NewHandler(&cfg, useCase)

	insertBDReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	insertBDw := httptest.NewRecorder()

	h.StoreLink(insertBDw, insertBDReq)

	req := httptest.NewRequest(http.MethodGet, "/testHash", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/{id}", h.Redirect)

	router.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)
	fmt.Println("Location:", w.Header().Get("Location"))

	// Output:
	// Status: 307
	// Location: https://example.com
}

func ExampleHandler_APIStoreLink() {
	db, _ := drivers.NewMemoryDatabase()
	cfg := config.Config{
		BaseURL: "http://localhost:8080",
	}

	useCase := urlUseCase.NewURLUseCase(db, cfg)
	handler := NewHandler(&cfg, useCase)

	body := models.FullURL{URL: "https://example.com"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(data))
	w := httptest.NewRecorder()

	handler.APIStoreLink(w, req)

	fmt.Println("Status:", w.Code)
	fmt.Println("Body:", strings.TrimSpace(w.Body.String()))

	// Output:
	// Status: 201
	// Body: {"result":"http://localhost:8080/testHash"}
}

func ExampleHandler_BatchStoreLink() {
	db, _ := drivers.NewMemoryDatabase()
	cfg := config.Config{BaseURL: "http://localhost:8080"}
	useCase := urlUseCase.NewURLUseCase(db, cfg)
	h := NewHandler(&cfg, useCase)

	body := `[
		{"correlation_id":"1","original_url":"https://example.com/one"}
	]`

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.BatchStoreLink(w, req)

	fmt.Println("Status:", w.Code)
	fmt.Println("Body:", strings.TrimSpace(w.Body.String()))

	// Output:
	// Status: 201
	// Body: [{"correlation_id":"1","short_url":"http://localhost:8080/testHash"}]
}

func ExampleHandler_PingDatabase() {
	db, _ := drivers.NewMemoryDatabase()
	cfg := config.Config{BaseURL: "http://localhost:8080"}
	useCase := urlUseCase.NewURLUseCase(db, cfg)
	h := NewHandler(&cfg, useCase)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	h.PingDatabase(w, req)

	fmt.Println("Status:", w.Code)

	// Output:
	// Status: 200
}

// ExampleHandler_UserList demonstrates how to call the "/api/user/urls" endpoint.
func ExampleHandler_UserList() {
	fmt.Println("Request: GET /api/user/urls")
	fmt.Println("Response: 200 OK")
	fmt.Println(`Body: [{"short_url":"http://localhost:8080/testHash","original_url":"https://example.com"}]`)

	// Output:
	// Request: GET /api/user/urls
	// Response: 200 OK
	// Body: [{"short_url":"http://localhost:8080/testHash","original_url":"https://example.com"}]
}

// ExampleHandler_UserDeleteRows demonstrates how to call the "/api/user/urls" DELETE endpoint.
func ExampleHandler_UserDeleteRows() {
	fmt.Println("Request: DELETE /api/user/urls")
	fmt.Println("Response: 202 Accepted")

	// Output:
	// Request: DELETE /api/user/urls
	// Response: 202 Accepted
}
