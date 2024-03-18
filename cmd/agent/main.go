package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"Vova4o/metrix/internal/app"
	"Vova4o/metrix/internal/config"
	"Vova4o/metrix/internal/logger"
)

func main() {
	// Open a file for logging
	LogfileAgent, err := logger.Logger(config.AgentLogFile)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	if err == nil {
		logrus.SetOutput(LogfileAgent)
	} else {
		logrus.Info("Failed to open log file, using default stderr output.")
	}

	defer LogfileAgent.Close()

	// Set the output destination of the standard logger
	log.SetOutput(LogfileAgent)

	fmt.Println("Went to main")

	ctx := context.Background()
	client := resty.New()
	err = app.NewAgent(ctx, client)
	if err != nil {
		log.Fatalf("Failed to start the agent: %v", err)
	}

}
