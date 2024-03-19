package app

import (
	"context"

	"Vova4o/metrix/internal/clientmetrics"
	"Vova4o/metrix/internal/logger"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// NewAgent creates and starts a new Metrics agent.
// It never returns, running a continuous loop to collect and send metrics.
//
// Returns:
//
//	error: an error occurred while creating or running the agent.
func NewAgent(ctx context.Context, client *resty.Client, logger *logger.FileLogger) error {
	// Add a middleware logger
	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		logger.Logger.WithFields(logrus.Fields{
			"url": request.URL,
		}).Info("Sending request")

		return nil
	})

	client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		logger.Logger.WithFields(logrus.Fields{
			"status": response.StatusCode(),
			"body":   response.String(),
		}).Info("Received response")

		return nil
	})

	metrics := clientmetrics.NewMetrics(client) // Create new Metrics

	runMetricsLoop(ctx, metrics, logger)

	return nil
}

func runMetricsLoop(ctx context.Context, metrics *clientmetrics.Metrics, logger *logger.FileLogger) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-metrics.PollTicker.C:
			if err := metrics.PollMetrics(); err != nil {
				logger.Logger.WithError(err).Error("Failed to poll metrics")
			}
		case <-metrics.ReportTicker.C:
			if err := metrics.ReportMetrics(metrics.BaseURL); err != nil {
				logger.Logger.WithError(err).Error("Failed to report metrics")
			}
		}
	}
}
