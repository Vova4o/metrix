package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMetricValue(t *testing.T) {
	tests := []struct {
		name           string
		metricType     string
		metricName     string
		expectedStatus int
		expectedBody   string
	}{
		{"Gauge Test", "gauge", "test", http.StatusOK, "123.45"},
		{"Counter Test", "counter", "test", http.StatusOK, "678"},
		{"Invalid Metric Type", "wrong", "test", http.StatusNotFound, "invalid metric type: wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router
			r := gin.Default()

			// Create a mock storager with some metrics
			s := &mockStorager{
				gauges: map[string]float64{
					"test": 123.45,
				},
				counters: map[string]int64{
					"test": 678,
				},
			}

			// Register the handler
			r.GET("/metrics/:metricType/:metricName", MetricValue(s))

			// Create a test request
			req, _ := http.NewRequest("GET", "/metrics/"+tt.metricType+"/"+tt.metricName, nil)

			// Create a test response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			r.ServeHTTP(rr, req)

			// Check the response status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %v, got %v", tt.expectedStatus, rr.Code)
			}

			// Check the response body
			if gotBody := rr.Body.String(); gotBody != tt.expectedBody {
				t.Errorf("expected %v, got %v", tt.expectedBody, gotBody)
			}
		})
	}
}

func TestMetricValueJSON(t *testing.T) {
	testCases := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Gauge Test",
			body: map[string]interface{}{
				"type": "gauge",
				"id":   "test",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"test","type":"gauge","value":0}`,
		},
		{
			name: "Counter Test",
			body: map[string]interface{}{
				"type": "counter",
				"id":   "test",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"test","type":"counter","delta":0}`,
		},
		{
			name: "Invalid Metric Type",
			body: map[string]interface{}{
				"type": "wrong",
				"id":   "test",
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"invalid metric type: wrong"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router
			r := gin.Default()

			// Create a mock storager
			s := &mockStorager{
				gauges: map[string]float64{
					"test": 0,
				},
				counters: map[string]int64{
					"test": 0,
				},
			}

			// Register the handler
			r.POST("/metrics", MetricValueJSON(s))

			// Create a test request body
			jsonBody, _ := json.Marshal(tt.body)

			// Create a test request
			req, _ := http.NewRequest("POST", "/metrics", bytes.NewBuffer(jsonBody))

			// Set the Content-Type header to application/json
			req.Header.Set("Content-Type", "application/json")

			// Create a test response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			r.ServeHTTP(rr, req)

			// Check the response status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %v, got %v", tt.expectedStatus, rr.Code)
			}

			// Check the response body
			var gotBody, expectedBody map[string]interface{}
			json.Unmarshal([]byte(rr.Body.Bytes()), &gotBody)
			json.Unmarshal([]byte(tt.expectedBody), &expectedBody)

			if !reflect.DeepEqual(gotBody, expectedBody) {
				t.Errorf("expected %v, got %v", expectedBody, gotBody)
			}
		})
	}
}
