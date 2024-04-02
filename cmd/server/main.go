package main

import (
	"fmt"
	"log"

	appserver "Vova4o/metrix/internal/app/server"
	"Vova4o/metrix/internal/logger"
)

const (
	serverLogFile = "server.log"
)

func main() {
	// Open a file for logging
	err := logger.New(serverLogFile)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to initialize logger")
	}

	defer func() {
		if err := logger.Close(); err != nil {
			fmt.Printf("Failed to close log file: %v\n", err)
		}
	}()

	err = appserver.NewServer()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
