package appserver

import (
	"fmt"

	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
	mw "Vova4o/metrix/internal/middleware"
	"Vova4o/metrix/internal/serverflags"
	"Vova4o/metrix/internal/storage"

	"github.com/gin-gonic/gin"
)

func NewServer() error {
	// Set the mode to release
	gin.SetMode(gin.ReleaseMode)
	// Create a new router
	router := gin.Default()

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

	router.Use(mw.RequestLogger(), mw.GzipMiddleware())

	// Add the handlers to the router
	router.POST("/update/:metricType/:metricName/:metricValue", handlers.HandleUpdateText(memStorager))

	router.POST("/update/", handlers.HandleUpdateJSON(memStorager))

	router.GET("/", handlers.ShowMetrics(memStorager, tempFile))

	router.GET("/value/:metricType/:metricName", handlers.MetricValue(memStorager))
	router.POST("/value/", handlers.MetricValueJSON(memStorager))

	fmt.Printf("Starting server on %s\n", serverflags.GetServerAddress())

	// Start the server
	return router.Run(serverflags.GetServerAddress())
}
