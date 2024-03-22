package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"Vova4o/metrix/internal/storage"
)

//go:embed templates/*
var templates embed.FS

// ShowMetrics is an HTTP handler that shows all the metrics
func ShowMetrics(storage storage.StorageInterface, tempFile string) http.HandlerFunc {
	// Parse the template file
	tmpl, errFunc := ParseTemplate(tempFile)
	if errFunc != nil {
		return errFunc
	}

	// Return the actual handler function
	return func(w http.ResponseWriter, r *http.Request) {
		// Set the content type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Create a map of maps to hold the metrics
		data := map[string]interface{}{
			"GaugeMetrics":   GaugeMetricType{}.GetAll(storage),
			"CounterMetrics": CounterMetricType{}.GetAll(storage),
		}

		// Execute the template with the data
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ParseTemplate parses the template file and returns the parsed template
func ParseTemplate(tempFile string) (*template.Template, func(w http.ResponseWriter, r *http.Request)) {
	tmpl, err := template.ParseFS(templates, filepath.Join("templates", tempFile))
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return nil, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
	return tmpl, nil
}
