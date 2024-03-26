package mw

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func TestRequestLogger(t *testing.T) {
//     // Create a test hook for the logger
//     hook := test.NewLocal(logger.Log)

//     // Create a test handler
//     nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         w.Write([]byte("OK"))
//     })

//     // Wrap the handler with the RequestLogger middleware
//     handler := middleware.RequestLogger(nextHandler)

//     // Create a test request
//     req, err := http.NewRequest("GET", "/test", nil)
//     if err != nil {
//         t.Fatalf("Failed to create request: %v", err)
//     }

//     // Record the response
//     rr := httptest.NewRecorder()
//     handler.ServeHTTP(rr, req)

//     // Check the response status code
//     assert.Equal(t, http.StatusOK, rr.Code)

//     // Check the response body
//     assert.Equal(t, "OK", rr.Body.String())

//     // Check the log entry
//     assert.Equal(t, 1, len(hook.Entries))
//     entry := hook.LastEntry()
//     assert.Equal(t, "Handled request", entry.Message)
//     assert.Equal(t, "/test", entry.Data["path"])
//     assert.Equal(t, "GET", entry.Data["method"])
//     assert.Equal(t, http.StatusOK, entry.Data["status"])
//     assert.Contains(t, entry.Data["duration"], "ms")
//     assert.Equal(t, 2, entry.Data["size"])
// }

func TestGzipMiddleware(t *testing.T) {
	// Create a GzipMiddleware
	middleware := GzipMiddleware

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
		wantBody   string
	}{
		{
			name: "Test handler",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, world!"))
			}),
			wantStatus: http.StatusOK,
			wantBody:   "Hello, world!",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap the test handler with the middleware
			h := middleware(tt.handler)

			// Create a test request
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Accept-Encoding", "gzip")

			// Record the response
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
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

			if string(body) != tt.wantBody {
				t.Errorf("handler returned unexpected body: got %v want %v", string(body), tt.wantBody)
			}
		})
	}
}

func TestGzipMiddleware_NoGzip(t *testing.T) {
	// Create a GzipMiddleware
	middleware := GzipMiddleware

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
		wantBody   string
	}{
		{
			name: "Test handler",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, world!"))
			}),
			wantStatus: http.StatusOK,
			wantBody:   "Hello, world!",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap the test handler with the middleware
			h := middleware(tt.handler)

			// Create a test request without the "Accept-Encoding: gzip" header
			req := httptest.NewRequest("GET", "/", nil)

			// Record the response
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			// Check the Content-Encoding header
			if encoding := rr.Header().Get("Content-Encoding"); encoding == "gzip" {
				t.Errorf("handler returned wrong Content-Encoding header: got %v want %v", encoding, "")
			}

			// Check the response body
			body := rr.Body.String()
			if body != tt.wantBody {
				t.Errorf("handler returned unexpected body: got %v want %v", body, tt.wantBody)
			}
		})
	}
}
