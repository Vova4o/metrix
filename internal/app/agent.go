package app

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"Vova4o/metrix/internal/clientmetrics"
)

func NewAgent() error {
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

	m := clientmetrics.NewMetrics(client) // Create new Metrics

	// Start the main loop
	for {
		// Wait for the next tick
		select {
		// When the pollTicker ticks, we collect the metrics
		case <-m.PollTicker.C:
			m.PollMetrics()
		// When the reportTicker ticks, we send the metrics
		case <-m.ReportTicker.C:
			m.ReportMetrics(m.BaseURL)
		}
	}

}
