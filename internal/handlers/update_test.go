package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"Vova4o/metrix/internal/storage"
)

func Test_handleUpdate(t *testing.T) {
	tests := []struct {
		name        string
		metricType  string
		metricValue string
		wantStatus  int
	}{
		{
			name:        "Gauge metric",
			metricType:  "gauge",
			metricValue: "10.0",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "Counter metric",
			metricType:  "counter",
			metricValue: "10",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "Invalid metric value for counter",
			metricType:  "counter",
			metricValue: "invalid",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "Invalid metric value for gauge",
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

func TestHandleUpdateJSON(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Empty body",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			body:           `{"type":"gauge","name":"test","value":}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing name",
			body:           `{"type":"gauge","value":"10.0"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing value",
			body:           `{"type":"gauge","name":"test"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing type",
			body:           `{"name":"test","value":"10.0"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()

			// Create a request to pass to our handler
			req, err := http.NewRequest("POST", "", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
			rr := httptest.NewRecorder()
			handler := HandleUpdateJSON(storage)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
