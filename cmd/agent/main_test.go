package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type args struct {
	client     *resty.Client
	metricType string
	metricName string
	value      float64
}

func TestSendMetric(t *testing.T) {
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1: Valid gauge metric",
			args: args{
				client:     resty.New(),
				metricType: "gauge",
				metricName: "testMetric",
				value:      1.0,
			},
			wantErr: false,
		},
		{
			name: "Test case 2: Invalid metric type",
			args: args{
				client:     resty.New(),
				metricType: "invalid",
				metricName: "testMetric",
				value:      1.0,
			},
			wantErr: true,
		},
		// Add more test cases as needed
	}

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create a mock HTTP server
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            }))
            defer server.Close()

            // Set the mock server's URL as the base URL of the client
            tt.args.client.SetBaseURL(server.URL)

            // Call the function under test
            err := sendMetric(tt.args.client, tt.args.metricType, tt.args.metricName, tt.args.value)

            // Assert that there was no error
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
