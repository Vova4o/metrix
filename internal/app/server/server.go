package appserver

import (
	"fmt"
	"net/http"

	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
	mw "Vova4o/metrix/internal/middleware"
	"Vova4o/metrix/internal/serverflags"
	"Vova4o/metrix/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewServer() error {
	// Create a new router
	mux := chi.NewRouter()

	tempFile := "metrix.page.tmpl"

	// Create a new MemStorage
	memStorager := storage.NewMemStorage()

	if serverflags.GetFileStoragePath() != "" {
		fileStorage, err := storage.NewFileStorage(memStorager, serverflags.GetStoreInterval(), serverflags.GetFileStoragePath(), serverflags.GetRestore())
		if err != nil {
			err = fmt.Errorf("failed to create new file storage: %v", err)
			logger.Log.WithError(err).Error("Failed to create new file storage")
			return err
		}
		defer fileStorage.SaveToFile() // Save metrics to file on exit
	} else {
		fmt.Println("Not using file storage")
		logger.Log.Info("Not using file storage")
	}

	mux.Use(mw.RequestLogger)
	mux.Use(mw.GzipMiddleware)
	// mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Add the handlers to the router
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.HandleUpdateText(memStorager))
	mux.Post("/update/", handlers.HandleUpdateJSON(memStorager))

	mux.Get("/", handlers.ShowMetrics(memStorager, tempFile))

	mux.Get("/value/{metricType}/{metricName}", handlers.MetricValue(memStorager))
	mux.Post("/value/", handlers.MetricValueJSON(memStorager))

	fmt.Printf("Starting server on %s\n", serverflags.GetServerAddress())

	// Start the server
	return http.ListenAndServe(serverflags.GetServerAddress(), mux)
}
