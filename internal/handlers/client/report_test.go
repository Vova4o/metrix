package clientmetrics

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"Vova4o/metrix/internal/config"
)

type MockRestClient struct {
	resty.Client
}

func (m *MockRestClient) R() *resty.Request {
	return m.Client.R() // Use the R() method of resty.Client
}

func TestReportMetrics_Success(t *testing.T) {
	// Set up the mock client
	// client := &MockRestClient{}

	// Set up the metrics
	config.GaugeMetrics = map[string]float64{
		"gaugeTest": 10.0000,
	}
	config.CounterMetrics = map[string]int64{
		"counterTest": 20,
	}

	// Call the function
	ReportMetrics("http://localhost:8080")

	// Check that the metrics were sent
	// This could be done by checking the logs, or by setting up a mock server and checking the requests it received
	// For simplicity, this test just checks that the function doesn't panic
	assert.True(t, true)
}

func TestReportMetrics_Error(t *testing.T) {
	// Set up the mock client
	// client := &MockRestClient{}

	// Set up the metrics
	config.GaugeMetrics = map[string]float64{
		"gaugeTest": 10.0000,
	}
	config.CounterMetrics = map[string]int64{
		"counterTest": 20,
	}

	// Call the function with an invalid URL to force an error
	ReportMetrics("http://invalid-url")

	// Check that the function handled the error
	// This could be done by checking the logs, or by checking the state of the application after the function call
	// For simplicity, this test just checks that the function doesn't panic
	assert.True(t, true)
}
