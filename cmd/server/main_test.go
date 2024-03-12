package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestShowMetricsHandler(t *testing.T) {
	// Create a storage
	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
	}

	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "/metrics", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a HTTP handler
	handler := ShowMetrics(storage)

	// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the response body
	expected := "<html><body><h1>Gauge Metrics</h1><ul></ul><h1>Counter Metrics</h1><ul></ul></body></html>" // Expected response body
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body")
}

func Test_handleUpdate(t *testing.T) {
	// define an empty storage
	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
	}

	// Create a request to pass to our handler
	handler := handleUpdate(storage)

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
	// define an empty storage
	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
	}

	// Create a request to pass to our handler
	handler := handleUpdate(storage)

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
	// define an empty storage
	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
	}

	// Create a request to pass to our handler
	handler := handleUpdate(storage)

	// Create a new HTTP request with an invalid metric type
	req, err := http.NewRequest("POST", "", nil)
	// Check if there was an error creating the request
	assert.NoError(t, err)

	// Create a router context with the URL parameters
	rctx := chi.NewRouteContext()
	// Add the URL parameters with an invalid metric type
	rctx.URLParams.Add("metricType", "invalid")
	rctx.URLParams.Add("metricName", "test")
	rctx.URLParams.Add("metricValue", "10")
	// Add the router context to the request
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMetricValue(t *testing.T) {
	// define a storage with a gauge metric
	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
	}
	storage.SetGauge("test", 10.0)

	// Create a request to pass to our handler
	handler := MetricValue(storage)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "", nil)
	// Check if there was an error creating the request
	assert.NoError(t, err)

	// Create a router context with the URL parameters
	rctx := chi.NewRouteContext()
	// Add the URL parameters
	rctx.URLParams.Add("metricType", "gauge")
	rctx.URLParams.Add("metricName", "test")
	// Add the router context to the request
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expected := "10" // Expected response body
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body")
}

func TestMetricValue_notFound(t *testing.T) {
	// define an empty storage
	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
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
	rctx.URLParams.Add("metricType", "gauge")
	rctx.URLParams.Add("metricName", "test")
	// Add the router context to the request
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Call ServeHTTP method directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
