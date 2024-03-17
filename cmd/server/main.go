package main

import (
	"log"

	"Vova4o/metrix/internal/app"
	"Vova4o/metrix/internal/config"
	"Vova4o/metrix/internal/logger"
)

func main() {
	var err error
	// Open a file for logging
	config.LogfileServer, err = logger.Logger(config.ServerLogFile)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer config.LogfileServer.Close()

	log.SetOutput(config.LogfileServer)

	err = app.NewServer()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
