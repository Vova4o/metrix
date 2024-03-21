package main

import (
	"log"

	"Vova4o/metrix/internal/app"
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
	
	err = app.NewServer()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
