package appserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	mw "Vova4o/metrix/internal/MW"
	"Vova4o/metrix/internal/config"
	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
	"Vova4o/metrix/internal/serverflags"
	"Vova4o/metrix/internal/storage"
)

func NewServer() error {
	// Create a new router
	mux := chi.NewRouter()

	tempFile := "metrix.page.tmpl"

	// Create a new MemStorage
	memStorage := storage.NewMemStorage()

	// Create a new FileStorage
	fileStorage := storage.NewFileStorage(memStorage, serverflags.GetStoreInterval(), serverflags.GetFileStoragePath(), serverflags.GetRestore())
	defer fileStorage.SaveToFile() // Save metrics to file on exit

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
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.HandleUpdateText(fileStorage))
	mux.Post("/update/", handlers.HandleUpdateJSON(fileStorage))

	mux.Get("/", handlers.ShowMetrics(fileStorage, tempFile))

	mux.Get("/value/{metricType}/{metricName}", handlers.MetricValue(fileStorage))
	mux.Post("/value/", handlers.MetricValueJSON(fileStorage))

	fmt.Printf("Starting server on %s\n", serverflags.GetServerAddress())

	// Start the server
	return http.ListenAndServe(serverflags.GetServerAddress(), mux)
}
