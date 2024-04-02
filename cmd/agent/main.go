package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"

	appagent "Vova4o/metrix/internal/app/agent"
	"Vova4o/metrix/internal/logger"
)

const (
	agentLogFile = "agent.log"
)

func main() {
	err := logger.New(agentLogFile)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to initialize logger")
	}

	defer func() {
		if err := logger.Close(); err != nil {
			fmt.Printf("Failed to close log file: %v\n", err)
		}
	}()

	ctx := context.Background()
	client := resty.New()
	err = appagent.NewAgent(ctx, client)
	if err != nil {
		log.Fatalf("Failed to start the agent: %v", err)
	}
}
