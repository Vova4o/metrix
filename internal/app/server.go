package app

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"Vova4o/metrix/internal/config"
	serverflag "Vova4o/metrix/internal/flag"
	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/storage"
)

func NewServer() {

	// Create a new router
	mux := chi.NewRouter()

	// Create a new storage
	storage := &storage.MemStorage{
		GaugeMetrics:   sync.Map{},
		CounterMetrics: sync.Map{},
	}

	mux.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(config.LogfileServer, "", log.LstdFlags)}))
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Add the handlers to the router
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.HandleUpdate(storage))

	mux.Get("/", handlers.ShowMetrics(storage))

	mux.Get("/value/{metricType}/{metricName}", handlers.MetricValue(storage))

	fmt.Printf("Starting server on %s\n", *serverflag.ServerAddress)
	// Start the server
	http.ListenAndServe(*serverflag.ServerAddress, mux)

}
