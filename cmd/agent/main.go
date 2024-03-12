package main

import (
	"Vova4o/metrix/internal/clientflag"
	"Vova4o/metrix/internal/config"
	clientmetrics "Vova4o/metrix/internal/handlers/client"
	"Vova4o/metrix/internal/logger"
	"log"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	// Parse the flags
	clientflag.ParseFlags()

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

	//Start the cicle of collecting and sending metrics
	pollTicker := time.NewTicker(time.Duration(*clientflag.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(*clientflag.ReportInterval) * time.Second)
	baseURL := *clientflag.ServerAddress

	// Start the main loop
	for {
		// Wait for the next tick
		select {
		// When the pollTicker ticks, we collect the metrics
		case <-pollTicker.C:
			clientmetrics.PollMetrics()
			// When the reportTicker ticks, we send the metrics
		case <-reportTicker.C:
			clientmetrics.ReportMetrics(baseURL)
		}
	}
}
