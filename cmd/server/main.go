package main

import (
	"log"

	appserver "Vova4o/metrix/internal/app/server"
	"Vova4o/metrix/internal/config"
	"Vova4o/metrix/internal/logger"
)

func main() {
	// Open a file for logging
	logger, err := logger.NewLogger(config.ServerLogFile)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.CloseLogger()

	logger.SetOutput()

	err = appserver.NewServer()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
