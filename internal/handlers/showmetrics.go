package handlers

import (
    "fmt"
    "net/http"

    "Vova4o/metrix/internal/storage"
)

// ShowMetrics is an HTTP handler that shows all the metrics
func ShowMetrics(storage storage.StorageInterface) http.HandlerFunc {
    // Return the actual handler function
    return func(w http.ResponseWriter, r *http.Request) {
        // Get all the gauge metrics
        gaugeMetrics := storage.GetAllGauges()
        // Get all the counter metrics
        counterMetrics := storage.GetAllCounters()

        // Start the HTML response
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, "<html><body>")
        fmt.Fprint(w, "<h1>Gauge Metrics</h1><ul>")
        for key, value := range gaugeMetrics {
            fmt.Fprintf(w, "<li>%s: %.04f</li>", key, value)
        }
        fmt.Fprint(w, "</ul>")
        fmt.Fprint(w, "<h1>Counter Metrics</h1><ul>")
        for key, value := range counterMetrics {
            fmt.Fprintf(w, "<li>%s: %d</li>", key, int(value))
        }
        fmt.Fprint(w, "</ul></body></html>")
    }
}