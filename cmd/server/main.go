package main

import (
	"Vova4o/metrix/internal/config"
	serverflag "Vova4o/metrix/internal/flag"
	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
	"Vova4o/metrix/internal/methods"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Parse the flags
	serverflag.ParseFlags()

	// Open a file for logging
	LogfileServer, err := logger.Logger(config.ServerLogFile)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer LogfileServer.Close()

	log.SetOutput(LogfileServer)

	// Create a new router
	mux := chi.NewRouter()

	// Create a new storage
	storage := &methods.MemStorage{
		GaugeMetrics:   sync.Map{},
		CounterMetrics: sync.Map{},
	}

	mux.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(LogfileServer, "", log.LstdFlags)}))
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
