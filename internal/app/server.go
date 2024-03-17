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

func NewServer() error {

	// Create a new router
	mux := chi.NewRouter()

	// Create a new MemStorage
	memStorage := storage.NewMemStorage()

	mux.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(config.LogfileServer, "", log.LstdFlags)}))
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Add the handlers to the router
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.HandleUpdate(memStorage))
	// mux.Post("/update/json", handlers.HandleUpdateJSON(memStorage))

	mux.Get("/", handlers.ShowMetrics(memStorage))

	mux.Get("/value/{metricType}/{metricName}", handlers.MetricValue(memStorage))

	fmt.Printf("Starting server on %s\n", allflags.GetServerAddress())
	// Start the server
	return http.ListenAndServe(allflags.GetServerAddress(), mux)
}
