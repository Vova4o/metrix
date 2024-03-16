package app

import (
	"time"

	allflags "Vova4o/metrix/internal/flag"
	clientmetrics "Vova4o/metrix/internal/handlers/client"

	"github.com/go-resty/resty/v2"
)

func NewAgent() {
	// Create a new resty client
	client := resty.New()

	// Create a new MetricsCollector
	mc := &clientmetrics.MetricsCollector{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int),
	}

	// Create a new MetricsAgent
	ma := &clientmetrics.MetricsAgent{
		MetricsCollector: mc,
		Client:           client,
	}

	// Start the cycle of collecting and sending metrics
	pollTicker := time.NewTicker(time.Duration(allflags.GetPollInterval()) * time.Second)
	reportTicker := time.NewTicker(time.Duration(allflags.GetReportInterval()) * time.Second)
	baseURL := allflags.GetServerAddress()

	// Start the main loop
	for {
		// Wait for the next tick
		select {
		// When the pollTicker ticks, we collect the metrics
		case <-pollTicker.C:
			ma.PollMetrics()
		// When the reportTicker ticks, we send the metrics
		case <-reportTicker.C:
			ma.ReportMetrics(baseURL)
		}
	}
}
