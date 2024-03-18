package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"Vova4o/metrix/internal/storage"
)

// ShowMetrics is an HTTP handler that shows all the metrics
func ShowMetrics(storage storage.StorageInterface, tempFile string) http.HandlerFunc {
	
	// Parse the template file
	tmpl, err := template.ParseFiles(filepath.Join("../../templates", tempFile))
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	// Return the actual handler function
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all the gauge metrics
		gaugeMetrics := storage.GetAllGauges()
		// Get all the counter metrics
		counterMetrics := storage.GetAllCounters()

		// Create a map of maps to hold the metrics
		data := map[string]interface{}{
			"GaugeMetrics":   gaugeMetrics,
			"CounterMetrics": counterMetrics,
		}

		// Execute the template with the data
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
