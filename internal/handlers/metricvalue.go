package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"Vova4o/metrix/internal/storage"
)

func (g GaugeMetricType) GetValue(storage storage.StorageInterface, name string) (interface{}, bool) {
	return storage.GetGauge(name)
}

func (g GaugeMetricType) FormatValue(value interface{}) string {
	return strconv.FormatFloat(value.(float64), 'f', -1, 64)
}

func (c CounterMetricType) GetValue(storage storage.StorageInterface, name string) (interface{}, bool) {
	return storage.GetCounter(name)
}

func (c CounterMetricType) FormatValue(value interface{}) string {
	return fmt.Sprintf("%d", int(value.(int64)))
}

func MetricValue(storage storage.StorageInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

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

		value, exists := mt.GetValue(storage, metricName)
		if !exists {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mt.FormatValue(value))
	}
}
