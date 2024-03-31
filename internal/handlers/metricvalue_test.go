package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func (m *mockStorager) GetValue(metricType, metricName string) (float64, bool) {
	if metricType == "gauge" {
		value, ok := m.gauges[metricName]
		return value, ok
	} else if metricType == "counter" {
		value, ok := m.counters[metricName]
		return float64(value), ok
	}
	return 0, false
}

func TestMetricValue(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		metricType     string
		metricName     string
		expectedStatus int
		expectedBody   string
	}{
		{"Gauge Test", "gauge", "test", http.StatusOK, "123.45"},
		{"Counter Test", "counter", "test", http.StatusOK, "678"},
		{"Invalid Metric Type", "wrong", "test", http.StatusBadRequest, "Invalid metric type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router
			r := chi.NewRouter()

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
			r.Get("/metrics/{metricType}/{metricName}", MetricValue(s))

			// Create a test request
			req, err := http.NewRequest("GET", "/metrics/"+tt.metricType+"/"+tt.metricName, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a test response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			r.ServeHTTP(rr, req)

			// Check the response status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %v, got %v", tt.expectedStatus, rr.Code)
			}

			if gotBody := strings.TrimSpace(rr.Body.String()); gotBody != strings.TrimSpace(tt.expectedBody) {
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
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"Invalid metric type"`,
		},
		// Add more test cases as needed
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router
			r := chi.NewRouter()

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
			r.Post("/metrics", MetricValueJSON(s))

			// Create a test request body
			jsonBody, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}

			// Create a test request
			req, err := http.NewRequest("POST", "/metrics", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

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
			if tt.expectedStatus == http.StatusOK {
				var gotBody, expectedBody map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &gotBody); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(tt.expectedBody), &expectedBody); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(gotBody, expectedBody) {
					t.Errorf("expected %v, got %v", expectedBody, gotBody)
				}
			} else {
				gotBody := strings.Trim(string(rr.Body.Bytes()), "\n")
				if fmt.Sprintf("%q", gotBody) != tt.expectedBody {
					t.Errorf("expected %v, got %v", tt.expectedBody, fmt.Sprintf("%q", gotBody))
				}
			}
		})
	}
}
