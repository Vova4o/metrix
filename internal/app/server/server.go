package appserver

import (
	"database/sql"
	"fmt"

	flag "Vova4o/metrix/internal/flags/server"
	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
	mw "Vova4o/metrix/internal/middleware"
	"Vova4o/metrix/internal/storage"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func NewServer() error {
	// Set the mode to release
	gin.SetMode(gin.ReleaseMode)
	// Create a new router
	router := gin.Default()

	tempFile := "metrix.page.tmpl"

	// Create a new MemStorage
	memStorager := storage.NewMemory()

	if flag.FileStoragePath() != "" {
		fileStorage, err := storage.NewFile(memStorager, flag.StoreInterval(), flag.FileStoragePath(), flag.Restore())
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

	db, err := sql.Open("postgres", flag.DatabaseDSN())
	if err != nil {
		err = fmt.Errorf("failed to open database: %v", err)
		logger.Log.WithError(err).Error("Failed to open database")
		return err
	}

	router.Use(mw.RequestLogger())
	router.Use(mw.CompressGzip())
	router.Use(mw.DecompressGzip)
	// router.Use(mw.RequestLogger())

	// Add the handlers to the router
	router.POST("/update/:metricType/:metricName/:metricValue", handlers.HandleUpdateText(memStorager))

	router.POST("/update/", handlers.HandleUpdateJSON(memStorager))

	router.GET("/", handlers.ShowMetrics(memStorager, tempFile))

	router.GET("/value/:metricType/:metricName", handlers.MetricValue(memStorager))
	router.POST("/value/", handlers.MetricValueJSON(memStorager))

	router.GET("/ping", handlers.Ping(db))

	fmt.Printf("Starting server on %s\n", flag.ServerAddress())

	// Start the server
	return router.Run(flag.ServerAddress())
}
