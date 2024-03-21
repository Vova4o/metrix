package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"Vova4o/metrix/internal/storage"
)

func HandleUpdateText(storage storage.StorageInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		var mt MetricType
		switch metricType {
		case "gauge":
			mt = GaugeMetricType{}
		case "counter":
			mt = CounterMetricType{}
		default:
			log.Printf("Invalid metric type: %s", metricType)
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		value, err := mt.ParseValue(metricValue)
		if err != nil {
			logAndRespondError(w, err, "Invalid metric value", http.StatusBadRequest)
			return
		}

		mt.Store(storage, metricName, value)

		w.WriteHeader(http.StatusOK)
	}
}

func logAndRespondError(w http.ResponseWriter, err error, message string, code int) {
	log.Printf("%s: %s", message, err)
	http.Error(w, message, code)
}

func HandleUpdateJSON(storage storage.StorageInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics MetricsJSON
		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if metrics.ID == "" {
			http.Error(w, "Missing id", http.StatusBadRequest)
			return
		}

		var mt MetricType
		var value interface{}
		switch metrics.MType {
		case "gauge":
			mt = GaugeMetricType{}
			if metrics.Value != nil {
				value = *metrics.Value
			} else {
				http.Error(w, "Value is required for gauge type", http.StatusBadRequest)
				return
			}
		case "counter":
			mt = CounterMetricType{}
			if metrics.Delta != nil {
				value = *metrics.Delta
			}
		default:
			log.Printf("Invalid metric type: %s", metrics.MType)
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		mt.Store(storage, metrics.ID, value)

		// Get the latest value from the storage
		latestValue, ok := mt.GetValue(storage, metrics.ID)
		if !ok {
			http.Error(w, "Failed to get latest value", http.StatusInternalServerError)
			return
		}

		// Update the metrics value based on the type
		if metrics.MType == "counter" {
			if val, ok := latestValue.(int64); ok {
				metrics.Delta = &val
			} else {
				log.Printf("Expected *int64, got %T", latestValue)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else if metrics.MType == "gauge" {
			if val, ok := latestValue.(float64); ok {
				metrics.Value = &val
			} else {
				log.Printf("Expected *float64, got %T", latestValue)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)

		w.WriteHeader(http.StatusOK)
	}
}
