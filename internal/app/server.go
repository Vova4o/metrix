package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"Vova4o/metrix/internal/config"
	allflags "Vova4o/metrix/internal/flag"
	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/storage"
)

func NewServer() {

	// Create a new router
	mux := chi.NewRouter()

	// Create a new MemStorage
	memStorage := &storage.MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]float64),
	}

	mux.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(config.LogfileServer, "", log.LstdFlags)}))
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Add the handlers to the router
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.HandleUpdate(memStorage))

	mux.Get("/", handlers.ShowMetrics(memStorage))

	mux.Get("/value/{metricType}/{metricName}", handlers.MetricValue(memStorage))

	fmt.Printf("Starting server on %s\n", allflags.GetServerAddress())
	// Start the server
	http.ListenAndServe(allflags.GetServerAddress(), mux)

}
