package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"Vova4o/metrix/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestRequestLogger(t *testing.T) {
	// Create a test logger
	_ = logger.New("test.log")

	// Create a test hook for the logger
	hook := test.NewLocal(logger.Log)

	// Create a Gin engine for testing
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestLogger())
	engine.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	engine.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	if rr.Body.String() != "OK" {
		t.Errorf("Expected response body %q, but got %q", "OK", rr.Body.String())
	}

	// Check the log entry
	if len(hook.Entries) != 1 {
		t.Errorf("Expected 1 log entry, but got %d", len(hook.Entries))
	}
	entry := hook.LastEntry()
	if entry.Message != "Handled request" {
		t.Errorf("Expected log message %q, but got %q", "Handled request", entry.Message)
	}
	if entry.Data["path"] != "/test" {
		t.Errorf("Expected path %q, but got %q", "/test", entry.Data["path"])
	}
	if entry.Data["method"] != "GET" {
		t.Errorf("Expected method %q, but got %q", "GET", entry.Data["method"])
	}
	if entry.Data["status"] != http.StatusOK {
		t.Errorf("Expected status %d, but got %v", http.StatusOK, entry.Data["status"])
	}
	if _, ok := entry.Data["duration"]; !ok {
		t.Error("Duration key is missing in log entry")
	}
	if entry.Data["size"] != 2 {
		t.Errorf("Expected size %d, but got %v", 2, entry.Data["size"])
	}
}

func TestGzipMiddleware(t *testing.T) {
	// Create a Gin engine for testing
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(GzipMiddleware())
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	// Create a test request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	// Record the response
	rr := httptest.NewRecorder()
	engine.ServeHTTP(rr, req)

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

func TestGzipMiddleware_NoGzip(t *testing.T) {
	// Create a Gin engine for testing
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(GzipMiddleware())
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	// Create a test request without the "Accept-Encoding: gzip" header
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	engine.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the Content-Encoding header
	if encoding := rr.Header().Get("Content-Encoding"); encoding == "gzip" {
		t.Errorf("handler returned wrong Content-Encoding header: got %v want %v", encoding, "")
	}

	// Check the response body
	body := rr.Body.String()
	if body != "Hello, world!" {
		t.Errorf("handler returned unexpected body: got %v want %v", body, "Hello, world!")
	}
}
