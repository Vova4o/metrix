package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MetricValue(s Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		var mt Metricer
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

		value, exists := mt.GetValue(s, metricName)
		if !exists {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mt.FormatValue(value))
	}
}

func MetricValueJSON(s Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics MetricsJSON

		// Decode the JSON request body into the metrics struct
		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var mt Metricer
		switch metrics.MType {
		case "gauge":
			mt = GaugeMetricType{}
		case "counter":
			mt = CounterMetricType{}
		default:
			log.Printf("Invalid metric type: %s", metrics.MType)
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		value, exists := mt.GetValue(s, metrics.ID)
		if !exists {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if metrics.MType == "gauge" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    metrics.ID,
				"type":  metrics.MType,
				"value": value,
			})
			return
		}
		if metrics.MType == "counter" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    metrics.ID,
				"type":  metrics.MType,
				"delta": value,
			})
		}
	}
}
