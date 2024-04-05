package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
	// Create a mock Storager
	mockStorager := &mockStorager{}

	// Define the test cases
	testCases := []struct {
		name           string
		body           string
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Valid JSON",
			body:           `[{"id": "test", "type": "gauge", "value": 1.0}]`,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"id": "test", "type": "gauge", "value": 1.0},
		},
		{
			name:           "Valid single JSON for gauge",
			body:           `{"id": "test", "type": "gauge", "value": 1.0}`,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"id": "test", "type": "gauge", "value": 1.0},
		},
		{
			name:           "Valid JSON with multiple metrics",
			body:           `[{"id": "test1", "type": "gauge", "value": 1.0}, {"id": "test2", "type": "gauge", "value": 2.0}]`,
			expectedStatus: http.StatusOK,
			expectedBody:   []map[string]interface{}{{"id": "test1", "type": "gauge", "value": 1.0}, {"id": "test2", "type": "gauge", "value": 2.0}},
		},
		{
			name:           "Invalid JSON",
			body:           `{"id": "test", "type": "invalid", "value": 1.0]`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "Invalid JSON"},
		},
		// ... other test cases ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new gin context with the test case body
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))

			// Call the function with the mock Storager and gin context
			HandleUpdateJSON(mockStorager)(c)

			// Check the response status and body
			assert.Equal(t, tc.expectedStatus, w.Code)

			var responseBody interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			responseAsMap, ok := responseBody.(map[string]interface{})
			if ok {
				assert.Equal(t, tc.expectedBody, responseAsMap)
			} else {
				responseAsSlice, ok := responseBody.([]interface{})
				if !ok {
					t.Fatalf("Failed to assert response body type to []interface{}")
				}

				var response []map[string]interface{}
				for _, item := range responseAsSlice {
					itemAsMap, ok := item.(map[string]interface{})
					if !ok {
						t.Fatalf("Failed to assert response body item type to map[string]interface{}")
					}
					response = append(response, itemAsMap)
				}

				assert.Equal(t, tc.expectedBody, response)
			}
		})
	}
}
