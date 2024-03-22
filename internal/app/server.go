package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	mw "Vova4o/metrix/internal/MW"
	"Vova4o/metrix/internal/config"
	allflags "Vova4o/metrix/internal/flag"
	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
	"Vova4o/metrix/internal/storage"
)

func NewServer() error {
	// Create a new router
	mux := chi.NewRouter()

	tempFile := "metrix.page.tmpl"

	// Create a new MemStorage
	memStorage := storage.NewMemStorage()

	// Create a new FileLogger
	log, err := logger.NewLogger(config.ServerLogFile)
	if err != nil {
		return err
	}
	defer log.CloseLogger()

	mux.Use(mw.RequestLogger(log))
	mux.Use(mw.GzipMiddleware)
	// mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Add the handlers to the router
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.HandleUpdateText(memStorage))
	mux.Post("/update/", handlers.HandleUpdateJSON(memStorage))

	mux.Get("/", handlers.ShowMetrics(memStorage, tempFile))

	mux.Get("/value/{metricType}/{metricName}", handlers.MetricValue(memStorage))
	mux.Post("/value/", handlers.MetricValueJSON(memStorage))

	fmt.Printf("Starting server on %s\n", allflags.GetServerAddress())

	// Start the server
	return http.ListenAndServe(allflags.GetServerAddress(), mux)
}
