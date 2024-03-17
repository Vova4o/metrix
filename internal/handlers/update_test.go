package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"Vova4o/metrix/internal/storage"
)

func Test_handleUpdate(t *testing.T) {
	// define an empty storage
	storage := storage.NewMemStorage()

	// Create a request to pass to our handler
	handler := HandleUpdateText(storage)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "", nil)
	// Check if there was an error creating the request
	assert.NoError(t, err)

	// Create a router context with the URL parameters
	rctx := chi.NewRouteContext()
	// Add the URL parameters
	rctx.URLParams.Add("metricType", "gauge")
	rctx.URLParams.Add("metricName", "test")
	rctx.URLParams.Add("metricValue", "10.0")
	// Add the router context to the request
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_handleUpdate_counter(t *testing.T) {
	storage := storage.NewMemStorage()

	// Create a request to pass to our handler
	handler := HandleUpdateText(storage)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "", nil)
	// Check if there was an error creating the request
	assert.NoError(t, err)

	// Create a router context with the URL parameters
	rctx := chi.NewRouteContext()
	// Add the URL parameters
	rctx.URLParams.Add("metricType", "counter")
	rctx.URLParams.Add("metricName", "test")
	rctx.URLParams.Add("metricValue", "10")
	// Add the router context to the request
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_handleUpdate_error(t *testing.T) {
	tests := []struct {
		name        string
		metricType  string
		metricValue string
		wantStatus  int
	}{
		{
			name:        "Invalid metric value",
			metricType:  "counter",
			metricValue: "invalid",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "Invalid metric value",
			metricType:  "gauge",
			metricValue: "invalid",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "Invalid metric type",
			metricType:  "invalid",
			metricValue: "10",
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()

			// Create a request to pass to our handler
			handler := HandleUpdateText(storage)

			// Create a new HTTP request
			req, err := http.NewRequest("POST", "", nil)
			// Check if there was an error creating the request
			assert.NoError(t, err)

			// Create a router context with the URL parameters
			rctx := chi.NewRouteContext()
			// Add the URL parameters
			rctx.URLParams.Add("metricType", tt.metricType)
			rctx.URLParams.Add("metricName", "test")
			rctx.URLParams.Add("metricValue", tt.metricValue)
			// Add the router context to the request
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}
