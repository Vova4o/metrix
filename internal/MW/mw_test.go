package mw

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"Vova4o/metrix/internal/logger"
)

func TestRequestLogger(t *testing.T) {
	// Create a FileLogger
	fileLogger, _ := logger.NewLogger("test.log")

	// Create a RequestLogger middleware
	middleware := RequestLogger(fileLogger)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// Wrap the test handler with the middleware
	h := middleware(handler)

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)

	// Record the response
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGzipMiddleware(t *testing.T) {
	// Create a GzipMiddleware
	middleware := GzipMiddleware

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// Wrap the test handler with the middleware
	h := middleware(handler)

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	// Record the response
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the Content-Encoding header
	if encoding := rr.Header().Get("Content-Encoding"); encoding != "gzip" {
		t.Errorf("handler returned wrong Content-Encoding header: got %v want %v", encoding, "gzip")
	}

	// Check the response body
	reader, err := gzip.NewReader(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "Hello, world!" {
		t.Errorf("handler returned unexpected body: got %v want %v", string(body), "Hello, world!")
	}
}
