package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"Vova4o/metrix/internal/storage"
)

func NewMemStorage() *storage.MemStorage {
	return &storage.MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]float64),
	}
}

func TestShowMetricsHandler(t *testing.T) {
	// Create a storage and set some metrics
	storage := NewMemStorage()
	storage.SetGauge("gaugeTest", 10.0000)
	storage.SetCounter("counterTest", 20)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := ShowMetrics(storage)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body is what we expect
	expected := `<html><body><h1>Gauge Metrics</h1><ul><li>gaugeTest: 10.0000</li></ul><h1>Counter Metrics</h1><ul><li>counterTest: 20</li></ul></body></html>` // fill this with the expected response
	assert.Equal(t, expected, rr.Body.String())
}
