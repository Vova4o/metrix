package app

import (
	"time"

	"Vova4o/metrix/internal/clientmetrics"
	allflags "Vova4o/metrix/internal/flag"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

func NewAgent() {
	client := resty.New()

	// Add a middleware logger
	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		logrus.WithFields(logrus.Fields{
			"url": request.URL,
		}).Info("Sending request")

		return nil
	})

	client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		logrus.WithFields(logrus.Fields{
			"status": response.StatusCode(),
			"body":   response.String(),
		}).Info("Received response")

		return nil
	})

	m := clientmetrics.NewMetrics() // Create new Metrics

	ma := &clientmetrics.MetricsAgent{
		Metrics: m, // Use Metrics instead of Gauge and Counter
		Client:  client,
	}

	pollTicker := time.NewTicker(time.Duration(allflags.GetPollInterval()) * time.Second)
	reportTicker := time.NewTicker(time.Duration(allflags.GetReportInterval()) * time.Second)
	baseURL := allflags.GetServerAddress()

	// Start the main loop
	for {
		// Wait for the next tick
		select {
		// When the pollTicker ticks, we collect the metrics
		case <-pollTicker.C:
			ma.Metrics.PollMetrics()
		// When the reportTicker ticks, we send the metrics
		case <-reportTicker.C:
			ma.ReportMetrics(baseURL)
		}
	}

}

// package app

// import (
// 	"time"

// 	"Vova4o/metrix/internal/clientmetrics"
// 	allflags "Vova4o/metrix/internal/flag"

// 	"github.com/go-resty/resty/v2"
// )

// func NewAgent() {
// 	// Create a new resty client
// 	client := resty.New()

// 	// Create a new MetricsCollector
// 	mc := &clientmetrics.MetricsCollector{
// 		GaugeMetrics:   make(map[string]float64),
// 		CounterMetrics: make(map[string]float64),
// 	}

// 	// Create a new MetricsAgent
// 	ma := &clientmetrics.MetricsAgent{
// 		MetricsCollector: mc,
// 		Client:           client,
// 	}

// 	// Start the cycle of collecting and sending metrics
// 	pollTicker := time.NewTicker(time.Duration(allflags.GetPollInterval()) * time.Second)
// 	reportTicker := time.NewTicker(time.Duration(allflags.GetReportInterval()) * time.Second)
// 	baseURL := allflags.GetServerAddress()

// 	// Start the main loop
// 	for {
// 		// Wait for the next tick
// 		select {
// 		// When the pollTicker ticks, we collect the metrics
// 		case <-pollTicker.C:
// 			ma.PollMetrics()
// 		// When the reportTicker ticks, we send the metrics
// 		case <-reportTicker.C:
// 			ma.ReportMetrics(baseURL)
// 		}
// 	}
// }
