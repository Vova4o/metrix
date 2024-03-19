package main

import (
	"context"
	"log"

	"github.com/go-resty/resty/v2"

	"Vova4o/metrix/internal/app"
	"Vova4o/metrix/internal/config"
	"Vova4o/metrix/internal/logger"
)

func main() {
	// Open a file for logging
	logger, err := logger.NewFileLogger(config.AgentLogFile)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.CloseLogger()

	logger.SetOutput()

	ctx := context.Background()
	client := resty.New()
	err = app.NewAgent(ctx, client, logger)
	if err != nil {
		log.Fatalf("Failed to start the agent: %v", err)
	}
}
