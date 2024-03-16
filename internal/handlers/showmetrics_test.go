package handlers

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"Vova4o/metrix/internal/storage"
)

func TestShowMetricsHandler(t *testing.T) {
	// Create a storage and set some metrics
	storage := &storage.MemStorage{
		GaugeMetrics:   sync.Map{},
		CounterMetrics: sync.Map{},
	}
	storage.SetGauge("gaugeTest", 10.0000)
	storage.SetCounter("counterTest", 20)

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
	expected := "<html><body><h1>Gauge Metrics</h1><ul><li>gaugeTest: 10.0000</li></ul><h1>Counter Metrics</h1><ul><li>counterTest: 20</li></ul></body></html>" // Expected response body
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body")
}