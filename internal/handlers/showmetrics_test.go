package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/storage"
)


func TestShowMetricsHandler_Success(t *testing.T) {
	// Create a storage and set some metrics
	storage := storage.NewMemStorage()
	storage.SetGauge("gaugeTest", 10.0000)
	storage.SetCounter("counterTest", 20)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := handlers.ShowMetrics(storage)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body contains the expected metrics and their values
	body := rr.Body.String()
	assert.Contains(t, body, "<li>gaugeTest: 10.0000</li>")
	assert.Contains(t, body, "<li>counterTest: 20</li>")
}
