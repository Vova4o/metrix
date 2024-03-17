package clientmetrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
)

type MockRestClient struct {
	client *resty.Client
}

func (m *MockRestClient) R() *resty.Request {
	return m.client.R()
}

func TestSendMetric(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        if strings.HasPrefix(req.URL.Path, "/error") {
            rw.WriteHeader(http.StatusInternalServerError)
            return
        }
        rw.WriteHeader(http.StatusOK)
    }))
    defer server.Close()

    tests := []struct {
        name    string
        sender  MetricSender
        client  *MockRestClient
        url     string
        wantErr bool
    }{
        {
            name: "Test Case 1 - Valid Client and Metrics with TextSender",
            sender: &TextMetricSender{},
            client: &MockRestClient{
                client: resty.New(),
            },
            url:     server.URL,
            wantErr: false,
        },
        {
            name: "Test Case 2 - Valid Client and Metrics with JSONSender",
            sender: &JSONMetricSender{},
            client: &MockRestClient{
                client: resty.New(),
            },
            url:     server.URL,
            wantErr: false,
        },
        {
            name: "Test Case 3 - Nil Client with TextSender",
            sender: &TextMetricSender{},
            client: &MockRestClient{
                client: nil,
            },
            url:     server.URL,
            wantErr: true,
        },
        {
            name: "Test Case 4 - Invalid URL with JSONSender",
            sender: &JSONMetricSender{},
            client: &MockRestClient{
                client: resty.New(),
            },
            url:     "",
            wantErr: true,
        },
        {
            name: "Test Case 5 - Status Code 500 with TextSender",
            sender: &TextMetricSender{},
            client: &MockRestClient{
                client: resty.New(),
            },
            url:     server.URL + "/error",
            wantErr: true,
        },
        // Add more test cases as needed
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.sender.SendMetric(tt.client.client, "gauge", "RandomValue", "0.5", tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("SendMetric() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}