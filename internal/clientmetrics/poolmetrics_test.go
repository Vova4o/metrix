package clientmetrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestMetrics_PollMetrics(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]float64
		CounterMetrics map[string]int64
	}
	tests := []struct {
		name             string
		fields           fields
		wantCounterValue int64
	}{
		{
			name: "Test Case 1",
			fields: fields{
				GaugeMetrics:   map[string]float64{"RandomValue": 0.5},
				CounterMetrics: map[string]int64{"CounterValue": 3},
			},
			wantCounterValue: 3, // expected value after PollMetrics
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			m.PollMetrics()
			randomValue, ok := m.GaugeMetrics["RandomValue"]
			if !ok {
				t.Errorf("RandomValue is missing")
			} else if randomValue < 0 || randomValue > 1 {
				t.Errorf("RandomValue = %v, want a value between 0 and 1", randomValue)
			}
			counterValue, ok := m.CounterMetrics["CounterValue"]
			if !ok {
				t.Errorf("CounterValue is missing")
			} else if counterValue != tt.wantCounterValue {
				t.Errorf("CounterValue = %v, want %v", counterValue, tt.wantCounterValue)
			}
		})
	}
}

func TestMetricsAgent_ReportMetrics(t *testing.T) {
	type fields struct {
		Metrics *Metrics
		Client  *resty.Client
	}
	type args struct {
		baseURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Case 1 - Valid Metrics and Client",
			fields: fields{
				Metrics: &Metrics{
					GaugeMetrics:   map[string]float64{"RandomValue": 0.5},
					CounterMetrics: map[string]int64{"CounterValue": 3},
				},
				Client: resty.New(),
			},
			args: args{
				baseURL: "http://localhost:8080/metrics", // replace with your actual server URL
			},
			wantErr: false,
		},
		{
			name: "Test Case 2 - Nil Metrics",
			fields: fields{
				Metrics: nil,
				Client:  resty.New(),
			},
			args: args{
				baseURL: "http://localhost:8080/metrics", // replace with your actual server URL
			},
			wantErr: true,
		},
		{
			name: "Test Case 3 - Nil Client",
			fields: fields{
				Metrics: &Metrics{
					GaugeMetrics:   map[string]float64{"RandomValue": 0.5},
					CounterMetrics: map[string]int64{"CounterValue": 3},
				},
				Client: nil,
			},
			args: args{
				baseURL: "http://localhost:8080/metrics", // replace with your actual server URL
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ma := &MetricsAgent{
				Metrics: tt.fields.Metrics,
				Client:  tt.fields.Client,
			}
			if err := ma.ReportMetrics(tt.args.baseURL); (err != nil) != tt.wantErr {
				t.Errorf("MetricsAgent.ReportMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MockRestClient struct {
	client *resty.Client
}

func (m *MockRestClient) R() *resty.Request {
	return m.client.R()
}

func TestSendMetric(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tests := []struct {
		name    string
		client  RestClient
		wantErr bool
	}{
		{
			name: "Test Case 1 - Valid Client and Metrics",
			client: &MockRestClient{
				client: resty.New(),
			},
			wantErr: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendMetric(tt.client, "gauge", "RandomValue", "0.5", server.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
