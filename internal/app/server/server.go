package appserver

import (
	"errors"
	"fmt"
	"time"

	flag "Vova4o/metrix/internal/flags/server"
	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
	mw "Vova4o/metrix/internal/middleware"
	"Vova4o/metrix/internal/storage"

	"github.com/cenkalti/backoff/v4"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
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

	var db *storage.DbStorage
	dsn := flag.DatabaseDSN()
	if dsn != "" {
		var err error

		operation := func() error {
			db, err = storage.NewDBConnection(dsn)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					// Check if the error code is retryable
					if pgErr.Code == "40001" || pgErr.Code == "40P01" {
						return err
					}
				}
				return backoff.Permanent(err)
			}
			return nil
		}
		b := backoff.NewExponentialBackOff()
		b.InitialInterval = 1 * time.Second
		b.MaxInterval = 5 * time.Second
		b.MaxElapsedTime = 10 * time.Second
		err = backoff.Retry(operation, b)
		if err != nil {
			err = fmt.Errorf("failed to create new database: %v", err)
			logger.Log.WithError(err).Error("Failed to create new database")
			return err
		}
		defer db.DB.Close()
	} else if flag.FileStoragePath() != "" {
		fileStorage, err := storage.NewFile(memStorager, flag.StoreInterval(), flag.FileStoragePath(), flag.Restore())
		if err != nil {
			err = fmt.Errorf("failed to create new file storage: %v", err)
			logger.Log.WithError(err).Error("Failed to create new file storage")
			return err
		}
		defer fileStorage.SaveTo() // Save metrics to file on exit
	}

	router.Use(mw.RequestLogger())
	router.Use(mw.CompressGzip())
	router.Use(mw.DecompressGzip)

	// Add the handlers to the router
	router.POST("/update/:metricType/:metricName/:metricValue", handlers.HandleUpdateText(memStorager))

	router.POST("/update/", handlers.HandleUpdateJSON(memStorager))
	router.POST("/updates/", handlers.HandleUpdateJSON(memStorager))

	router.GET("/", handlers.ShowMetrics(memStorager, tempFile))

	router.GET("/value/:metricType/:metricName", handlers.MetricValue(memStorager))
	router.POST("/value/", handlers.MetricValueJSON(memStorager))

	router.GET("/ping", handlers.Ping(db.DB))

	fmt.Printf("Starting server on %s\n", flag.ServerAddress())

	// Start the server
	return router.Run(flag.ServerAddress())
}
