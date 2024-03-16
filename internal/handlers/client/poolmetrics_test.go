package clientmetrics

import (
	"testing"

	"Vova4o/metrix/internal/config"
)

func TestPollMetrics(t *testing.T) {
	// Save the initial state of the metrics
	initialGaugeMetrics := config.GaugeMetrics
	initialCounterMetrics := config.CounterMetrics

	// Call the function
	PollMetrics()

	// Check that the gauge metrics have been updated
	if len(config.GaugeMetrics) == 0 {
		t.Errorf("GaugeMetrics was not updated")
	}

	if config.GaugeMetrics["RandomValue"] == initialGaugeMetrics["RandomValue"] {
		t.Errorf("GaugeMetrics was not updated")
	}
	// Check that the counter metrics have been updated
	if config.CounterMetrics["PoolCount"] != initialCounterMetrics["PoolCount"] {
		t.Errorf("CounterMetrics was not updated")
	}
}
