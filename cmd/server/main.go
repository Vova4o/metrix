package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

// parseFlags parses the flags and sets the serverAddress variable
var ServerAddress = flag.String("a", "", "HTTP server address")

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
	SetCounter(key string, value float64)
	//GetCounter returns the value of a counter
	GetCounter(key string) (float64, bool)
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

// SetCounter increments the value of a counter
func (m *MemStorage) SetCounter(key string, value float64) {
	actual, loaded := m.counterMetrics.LoadOrStore(key, value)
	if loaded {
		newValue := actual.(float64) + value
		m.counterMetrics.Store(key, newValue)
	}
}

// GetCounter returns the value of a counter
func (m *MemStorage) GetCounter(key string) (float64, bool) {
	value, exists := m.counterMetrics.Load(key)
	if exists {
		return value.(float64), exists
	}
	return 0, exists
}

// Delete removes a metric from the storage
func (m *MemStorage) Delete(key string) {
	m.gaugeMetrics.Delete(key)
	m.counterMetrics.Delete(key)
}

// handleUpdate is an HTTP handler that updates a metric
func handleUpdate(storage *MemStorage) http.HandlerFunc {
	// Return the actual handler function
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the metric type, name and value from the URL parameters
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		// Check if the metric type is valid
		switch metricType {
		case "gauge":
			// Parse the metric value as a float
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				log.Printf("Invalid metric value gauge: %s", metricValue)
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			// store the value in the storage
			storage.SetGauge(metricName, value)
		case "counter":
			// Parse the metric value as a float
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				log.Printf("Invalid metric value counter: %s", metricValue)
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			// store the value in the storage
			storage.SetCounter(metricName, value)
		default:
			log.Printf("Invalid metric type: %s", metricType)
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
		}
		// Set the status code to 200 OK
		w.WriteHeader(http.StatusOK)
	}
}

// ShowMetrics is an HTTP handler that shows all the metrics
func ShowMetrics(storage *MemStorage) http.HandlerFunc {
	// Return the actual handler function
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all the gauge metrics
		gaugeMetrics := map[string]float64{}
		// Range over all the gauge metrics and add them to the map
		storage.gaugeMetrics.Range(func(key, value interface{}) bool {
			// Add the metric to the map
			gaugeMetrics[key.(string)] = value.(float64)
			return true
		})

		// Get all the counter metrics
		counterMetrics := map[string]float64{}
		// Range over all the counter metrics and add them to the map
		storage.counterMetrics.Range(func(key, value interface{}) bool {
			// Add the metric to the map
			counterMetrics[key.(string)] = value.(float64)
			return true
		})

		// Start the HTML response
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<html><body>")
		fmt.Fprint(w, "<h1>Gauge Metrics</h1><ul>")
		for key, value := range gaugeMetrics {
			// Format the float with no trailing zeros
			formattedValue := strconv.FormatFloat(value, 'g', -1, 64)
			fmt.Fprintf(w, "<li>%s: %s</li>", key, formattedValue)
		}
		fmt.Fprint(w, "</ul>")
		fmt.Fprint(w, "<h1>Counter Metrics</h1><ul>")
		for key, value := range counterMetrics {
			fmt.Fprintf(w, "<li>%s: %d</li>", key, int(value))
		}
		fmt.Fprint(w, "</ul></body></html>")
	}
}

// MetricValue is an HTTP handler that returns the value of a metric
func MetricValue(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the metric type and name from the URL parameters
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		// Check if the metric type is valid
		switch metricType {
		// If the metric type is gauge
		case "gauge":
			// if value exists, set it to the value of the metric
			value, exists := storage.GetGauge(metricName)
			if !exists {
				http.Error(w, "Metric not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			// Format the float with no trailing zeros
			formattedValue := strconv.FormatFloat(value, 'f', -1, 64)
			fmt.Fprint(w, formattedValue)
			// If the metric type is counter
		case "counter":
			// if value exists, set it to the value of the metric
			value, exists := storage.GetCounter(metricName)
			if !exists {
				http.Error(w, "Metric not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%d", int(value))
		default:
			log.Printf("Invalid metric type: %s", metricType)
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
		}
	}
}

func main() {
	// Parse the flags
	parseFlags()

	// Creating logger, at some point i was done looking for mistakes manualy
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Create a new router
	mux := chi.NewRouter()

	// Create a new storage
	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
	}
	// Add the handlers to the router
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", handleUpdate(storage))

	mux.Get("/", ShowMetrics(storage))

	mux.Get("/value/{metricType}/{metricName}", MetricValue(storage))

	fmt.Printf("Starting server on %s\n", *ServerAddress)
	// Start the server
	http.ListenAndServe(*ServerAddress, mux)
}
