package main

import (
	"net/http"
	"strconv"
	"strings"
)

// MemStorage is a simple in-memory storage for metrics
type MemStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

// Metric is a generic metric type that can be used for any type of metric
// It is used to serialize and deserialize metrics to and from JSON
type Metric struct {
	Type  string  `json:"type"`
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

// StorageInterface is an interface for storage backends
type StorageInterface interface {
	//SetGauge sets the value of a gauge
	SetGauge(key string, value float64)
	//GetGauge returns the value of a gauge
	GetGauge(key string) (float64, bool)
	//SetCounter sets the value of a counter
	SetCounter(key string, value int64)
	//GetCounter returns the value of a counter
	GetCounter(key string) (int64, bool)
	//Delete removes a metric from the storage
	Delete(key string)
}

// SetGauge sets the value of a gauge
func (m *MemStorage) SetGauge(key string, value float64) {
	m.gaugeMetrics[key] = value
}

// GetGauge returns the value of a gauge
func (m *MemStorage) GetGauge(key string) (float64, bool) {
	value, exists := m.gaugeMetrics[key]
	return value, exists
}

// SetCounter sets the value of a counter
func (m *MemStorage) SetCounter(key string, value int64) {
	m.counterMetrics[key] += value
}

// GetCounter returns the value of a counter
func (m *MemStorage) GetCounter(key string) (int64, bool) {
	value, exists := m.counterMetrics[key]
	return value, exists
}

// Delete removes a metric from the storage
func (m *MemStorage) Delete(key string) {
	delete(m.gaugeMetrics, key)
	delete(m.counterMetrics, key)
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a new MemStorage that will be used to store the metrics
	// we create new storage inside the main function to make it unique for each instance of the server
	storage := &MemStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}

	// Register a handler for the /metrics endpoint
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			// If it's not, return a 405 Method Not Allowed error
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Split the URL path into parts
		parts := strings.Split(r.URL.Path, "/")
		// Check if the URL path has the correct format
		if len(parts) != 5 {
			// If it doesn't, return a 400 Bad Request error
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		// Extract the metric type, name, and value from the URL path
		metricType, metricName, metricValue := parts[2], parts[3], parts[4]

		// switch statement to handle different metric types
		switch metricType {
			// If the metric type is "gauge", parse the metric value as a float64
		case "gauge":
			//ParseFloat converts the string s to a floating-point number.
			value, err := strconv.ParseFloat(metricValue, 64)
			//Check if there is an error
			if err != nil {
				// If there is, return a 400 Bad Request error
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			// Set the value of the gauge metric in the storage
			storage.SetGauge(metricName, value)
			// If the metric type is "counter", parse the metric value as an int64
		case "counter":
			//ParseInt interprets a string.
			value, err := strconv.ParseInt(metricValue, 10, 64)
			//Check if there is an error
			if err != nil {
				// If there is, return a 400 Bad Request error
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			// Set the value of the counter metric in the storage
			storage.SetCounter(metricName, value)
			// If the metric type is neither "gauge" nor "counter", return a 400 Bad Request error
		default:
			//Send an error response and status code 400
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		// Return a 200 OK status code
		w.WriteHeader(http.StatusOK)
	})

	//Start the server on port 8080
	http.ListenAndServe("localhost:8080", mux)
}
