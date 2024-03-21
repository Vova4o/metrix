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
	_, err := logger.NewLogger(config.AgentLogFile)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Log.CloseLogger()

	logger.Log.SetOutput()

	ctx := context.Background()
	client := resty.New()
	err = app.NewAgent(ctx, client)
	if err != nil {
		log.Fatalf("Failed to start the agent: %v", err)
	}
}
