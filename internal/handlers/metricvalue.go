package handlers

import (
	"Vova4o/metrix/internal/methods"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)


// MetricValue is an HTTP handler that returns the value of a metric
func MetricValue(storage *methods.MemStorage) http.HandlerFunc {
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
