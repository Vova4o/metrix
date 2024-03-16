package clientmetrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendMetric(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a mock RestClient
	client := &MockRestClient{Client: *resty.New()}

	// Test the SendMetric function
	err := SendMetric(client, "gauge", "test", 10.0, server.URL)
	assert.NoError(t, err)

	// Test the SendMetric function with an invalid URL
	err = SendMetric(client, "gauge", "test", 10.0, "invalid-url")
	assert.Error(t, err)
}
