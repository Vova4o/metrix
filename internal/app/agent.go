package app

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"Vova4o/metrix/internal/clientmetrics"
)

// NewAgent creates and starts a new Metrics agent.
// It never returns, running a continuous loop to collect and send metrics.
//
// Returns:
//
//	error: an error occurred while creating or running the agent.
func NewAgent(ctx context.Context, client *resty.Client) error {

	fmt.Println("Hit the NewAgent function")

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

	metrics := clientmetrics.NewMetrics(client) // Create new Metrics

	runMetricsLoop(ctx, metrics)

	return nil
}

func runMetricsLoop(ctx context.Context, metrics *clientmetrics.Metrics) {
	fmt.Println("Hit the runMetricsLoop function")
	for {
		select {
		case <-ctx.Done():
			return
		case <-metrics.PollTicker.C:
			if err := metrics.PollMetrics(); err != nil {
				logrus.WithError(err).Error("Failed to poll metrics")
			}
		case <-metrics.ReportTicker.C:
			if err := metrics.ReportMetrics(metrics.BaseURL); err != nil {
				logrus.WithError(err).Error("Failed to report metrics")
			}
		}
	}
}
