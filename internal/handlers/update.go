package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"Vova4o/metrix/internal/storage"
)

type MetricUpdate struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (g GaugeMetricType) ParseValue(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}

func (g GaugeMetricType) Store(storage storage.StorageInterface, name string, value interface{}) {
	storage.SetGauge(name, value.(float64))
}

func (c CounterMetricType) ParseValue(value string) (interface{}, error) {
	return strconv.ParseInt(value, 10, 64)
}

func (c CounterMetricType) Store(storage storage.StorageInterface, name string, value interface{}) {
	storage.SetCounter(name, value.(int64))
}

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
		var update MetricUpdate
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var mt MetricType
		switch update.Type {
		case "gauge":
			mt = GaugeMetricType{}
		case "counter":
			mt = CounterMetricType{}
		default:
			log.Printf("Invalid metric type: %s", update.Type)
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		value, err := mt.ParseValue(update.Value)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}

		mt.Store(storage, update.Name, value)

		metrixName := "counter"
		// Get the latest value from the storage
		latestValue, ok := mt.GetValue(storage, metrixName)
		if !ok {
			http.Error(w, "Failed to get latest value", http.StatusInternalServerError)
			return
		}

		// Write the latest value to the response body
		response := struct {
			LatestValue interface{} `json:"delta"`
		}{
			LatestValue: latestValue,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		w.WriteHeader(http.StatusOK)
	}
}
