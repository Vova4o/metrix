package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"Vova4o/metrix/internal/storage"
)

func TestMetricValueUpdate(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		metricName string
		setValue   float64
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Test gauge metric",
			metricType: "gauge",
			metricName: "test",
			setValue:   10.0,
			wantStatus: http.StatusOK,
			wantBody:   "10",
		},
		{
			name:       "Test counter metric",
			metricType: "counter",
			metricName: "test",
			setValue:   10.0,
			wantStatus: http.StatusOK,
			wantBody:   "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// define a storage with a metric
			storage := storage.NewMemStorage()

			if tt.metricType == "gauge" {
				storage.SetGauge(tt.metricName, tt.setValue)
			} else if tt.metricType == "counter" {
				storage.SetCounter(tt.metricName, int64(tt.setValue))
			}

			// Create a request to pass to our handler
			handler := MetricValue(storage)

			// Create a new HTTP request
			req, err := http.NewRequest("GET", "", nil)
			// Check if there was an error creating the request
			assert.NoError(t, err)

			// Create a router context with the URL parameters
			rctx := chi.NewRouteContext()
			// Add the URL parameters
			rctx.URLParams.Add("metricType", tt.metricType)
			rctx.URLParams.Add("metricName", tt.metricName)
			// Add the router context to the request
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.wantStatus, rr.Code)

			// Check the response body
			assert.Equal(t, tt.wantBody, rr.Body.String(), "handler returned unexpected body")
		})
	}
}

func TestMetricValueNotFound(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		metricName string
		wantStatus int
	}{
		{
			name:       "Test not found gauge",
			metricType: "gauge",
			metricName: "test",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Test not found counter",
			metricType: "counter",
			metricName: "test",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Test not found default",
			metricType: "unknown",
			metricName: "test",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// define an empty storage
			storage := storage.NewMemStorage()

			// Create a request to pass to our handler
			handler := MetricValue(storage)

			// Create a new HTTP request
			req, err := http.NewRequest("GET", "", nil)
			// Check if there was an error creating the request
			assert.NoError(t, err)

			// Create a router context with the URL parameters
			rctx := chi.NewRouteContext()
			// Add the URL parameters
			rctx.URLParams.Add("metricType", tt.metricType)
			rctx.URLParams.Add("metricName", tt.metricName)
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

func TestMetricValueJson(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		metricName string
		setValue   float64
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Test gauge metric",
			metricType: "gauge",
			metricName: "test",
			setValue:   10.0,
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"test","type":"gauge","value":10}`,
		},
		{
			name:       "Test counter metric",
			metricType: "counter",
			metricName: "test",
			setValue:   10.0,
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"test","type":"counter","delta":10}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// define a storage with a metric
			storage := storage.NewMemStorage()

			if tt.metricType == "gauge" {
				storage.SetGauge(tt.metricName, tt.setValue)
			} else if tt.metricType == "counter" {
				storage.SetCounter(tt.metricName, int64(tt.setValue))
			}

			// Create a MetricsJSON object
			metrics := MetricsJSON{
				ID:    tt.metricName,
				MType: tt.metricType,
			}

			// Marshal the MetricsJSON object to a JSON string
			requestBody, err := json.Marshal(metrics)
			assert.NoError(t, err)

			// Create a request to pass to our handler
			handler := MetricValueJSON(storage)

			// Create a new HTTP request
			req, err := http.NewRequest("POST", "", bytes.NewBuffer(requestBody))
			// Check if there was an error creating the request
			assert.NoError(t, err)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.wantStatus, rr.Code)

			// Unmarshal the expected body into a map
			var expectedBody map[string]interface{}
			err = json.Unmarshal([]byte(tt.wantBody), &expectedBody)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			// Unmarshal the actual body into a map
			var actualBody map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			// Compare the maps
			assert.Equal(t, expectedBody, actualBody, "handler returned unexpected body")
		})
	}
}
