package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// MemStorage is a simple in-memory storage for metrics
// It uses two sync.Map to store gauge and counter metrics
// I had to change the type of the map to avoid concurrent map writes
type MemStorage struct {
	gaugeMetrics   sync.Map
	counterMetrics sync.Map
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
	m.gaugeMetrics.Store(key, value)
}

// GetGauge returns the value of a gauge
func (m *MemStorage) GetGauge(key string) (float64, bool) {
	value, exists := m.gaugeMetrics.Load(key)
	if exists {
		return value.(float64), exists
	}
	return 0, exists
}

// SetCounter sets the value of a counter
func (m *MemStorage) SetCounter(key string, value int64) {
	actual, loaded := m.counterMetrics.LoadOrStore(key, value)
	if loaded {
		m.counterMetrics.Store(key, actual.(int64)+value)
	}
}

// GetCounter returns the value of a counter
func (m *MemStorage) GetCounter(key string) (int64, bool) {
	value, exists := m.counterMetrics.Load(key)
	if exists {
		return value.(int64), exists
	}
	return 0, exists
}

// Delete removes a metric from the storage
func (m *MemStorage) Delete(key string) {
	m.gaugeMetrics.Delete(key)
	m.counterMetrics.Delete(key)
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a new MemStorage that will be used to store the metrics
	// we create new storage inside the main function to make it unique for each instance of the server
	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
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
			// Originaly was 400, but tests ask for 404...
			http.Error(w, "Invalid URL format", http.StatusNotFound)
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
				// Originaly was 400, but tests ask for 404...
				http.Error(w, "Invalid metric value", http.StatusNotFound)
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
