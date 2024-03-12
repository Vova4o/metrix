package handlers

import (
	"Vova4o/metrix/internal/methods"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// ShowMetrics is an HTTP handler that shows all the metrics
func ShowMetrics(storage *methods.MemStorage) http.HandlerFunc {
	// Return the actual handler function
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all the gauge metrics
		gaugeMetrics := getMetrics(&storage.GaugeMetrics)
		// Get all the counter metrics
		counterMetrics := getMetrics(&storage.CounterMetrics)

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
func getMetrics(metrics *sync.Map) map[string]float64 {
	result := map[string]float64{}
	metrics.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(float64)
		return true
	})
	return result
}
