package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

// MockClient is a mock implementation of the RestClient interface
type MockClient struct {
	client *resty.Client
}

// R returns a new resty request
func (m *MockClient) R() *resty.Request {
	return m.client.R()
}

// TestPollMetrics tests the pollMetrics function
func TestSendMetric(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a mock HTTP client
	client := &MockClient{client: resty.New()}

	// Test sendMetric function
	err := sendMetric(client, "gauge", "testMetric", 1.0, server.URL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test sendMetric function with server error
	serverError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverError.Close()

	err = sendMetric(client, "gauge", "testMetric", 1.0, serverError.URL)
	if err == nil {
		t.Errorf("Expected error, got nil")
	} else {
		expectedErrorMessage := "server returned non-OK status for gauge metric testMetric: 500 Internal Server Error"
		if err.Error() != expectedErrorMessage {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
		}
	}
}
