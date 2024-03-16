package app

import (
	"time"

	allflags "Vova4o/metrix/internal/flag"
	clientmetrics "Vova4o/metrix/internal/handlers/client"
)

func NewAgent() {

	//Start the cicle of collecting and sending metrics
	pollTicker := time.NewTicker(time.Duration(allflags.GetPollInterval()) * time.Second)
	reportTicker := time.NewTicker(time.Duration(allflags.GetReportInterval()) * time.Second)
	baseURL := allflags.GetServerAddress()

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
