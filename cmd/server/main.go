package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
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

func handleUpdate(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			storage.SetGauge(metricName, value)
		case "counter":
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			storage.SetCounter(metricName, value)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ShowMetrics(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all the gauge metrics
		gaugeMetrics := map[string]float64{}
		storage.gaugeMetrics.Range(func(key, value interface{}) bool {
			gaugeMetrics[key.(string)] = value.(float64)
			return true
		})

		// Get all the counter metrics
		counterMetrics := map[string]int64{}
		storage.counterMetrics.Range(func(key, value interface{}) bool {
			counterMetrics[key.(string)] = value.(int64)
			return true
		})

		// Start the HTML response
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<html><body>")
		fmt.Fprint(w, "<h1>Gauge Metrics</h1><ul>")
		for key, value := range gaugeMetrics {
			fmt.Fprintf(w, "<li>%s: %f</li>", key, value)
		}
		fmt.Fprint(w, "</ul><h1>Counter Metrics</h1><ul>")
		for key, value := range counterMetrics {
			fmt.Fprintf(w, "<li>%s: %d</li>", key, value)
		}
		fmt.Fprint(w, "</ul></body></html>")
	}
}

func MetricValue(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		switch metricType {
		case "gauge":
			value, exists := storage.GetGauge(metricName)
			if !exists {
				http.Error(w, "Metric not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%f", value)
		case "counter":
			value, exists := storage.GetCounter(metricName)
			if !exists {
				http.Error(w, "Metric not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%d", value)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
		}
	}
}

func main() {
	parseFlags()

	mux := chi.NewRouter()

	storage := &MemStorage{
		gaugeMetrics:   sync.Map{},
		counterMetrics: sync.Map{},
	}

	mux.Post("/update/{metricType}/{metricName}/{metricValue}", handleUpdate(storage))

	mux.Get("/", ShowMetrics(storage))

	mux.Get("/value/{metricType}/{metricName}", MetricValue(storage))

	http.ListenAndServe(*serverAddress, mux)
}
