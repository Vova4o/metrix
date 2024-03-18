package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"Vova4o/metrix/internal/storage"
)

// HandleUpdate is an HTTP handler that updates a metric
func HandleUpdateText(storage storage.StorageInterface) http.HandlerFunc {
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
				logAndRespondError(w, err, "Invalid metric value", http.StatusBadRequest)
				return
			}
			// store the value in the storage
			storage.SetGauge(metricName, value)
		case "counter":
			// Parse the metric value as a float
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				logAndRespondError(w, err, "Invalid metric value", http.StatusBadRequest)
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

func logAndRespondError(w http.ResponseWriter, err error, message string, code int) {
	log.Printf("%s: %s", message, err)
	http.Error(w, message, code)
}

type Metric struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func HandleUpdateJSON(storage storage.StorageInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric Metric
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&metric); err != nil {
			logAndRespondError(w, err, "Failed to decode JSON", http.StatusBadRequest)
			return
		}
		switch metric.Type {
		case "gauge":
			value, err := strconv.ParseFloat(metric.Value, 64)
			if err != nil {
				logAndRespondError(w, err, "Invalid metric value", http.StatusBadRequest)
				return
			}
			storage.SetGauge(metric.Name, value)
		case "counter":
			value, err := strconv.ParseInt(metric.Value, 10, 64)
			if err != nil {
				logAndRespondError(w, err, "Invalid metric value", http.StatusBadRequest)
				return
			}
			storage.SetCounter(metric.Name, value)
		default:
			log.Printf("Invalid metric type: %s", metric.Type)
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
	}
}
