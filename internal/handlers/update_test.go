package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleUpdateText(t *testing.T) {
	// Create a mock Gin context
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Set the request parameters
	c.Params = gin.Params{
		{Key: "metricType", Value: "gauge"},
		{Key: "metricName", Value: "testMetric"},
		{Key: "metricValue", Value: "10"},
	}

	// Create a mock storager
	s := &mockStorager{}

	// Call the handler function
	HandleUpdateText(s)(c)

	// Check the response status code
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, c.Writer.Status())
	}

	// Add more assertions if needed
}

func TestHandleUpdateJSON(t *testing.T) {
	tests := []struct {
		name           string
		metrics        MetricsJSON
		expectedStatus int
	}{
		{
			name: "valid gauge metric",
			metrics: MetricsJSON{
				ID:    "test",
				MType: "gauge",
				Value: pointerToFloat64(1.23),
			},
			expectedStatus: http.StatusOK,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock Gin context
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Convert the MetricsJSON object to JSON
			body, err := json.Marshal(tt.metrics)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Set the request body
			c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Create a mock storager
			s := &mockStorager{}

			// Call the handler function
			HandleUpdateJSON(s)(c)

			// Check the response status code
			if c.Writer.Status() != tt.expectedStatus {
				t.Errorf("expected %v, got %v", tt.expectedStatus, c.Writer.Status())
			}
		})
	}
}

func pointerToFloat64(f float64) *float64 {
	return &f
}
