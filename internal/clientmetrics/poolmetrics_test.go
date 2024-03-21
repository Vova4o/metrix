package clientmetrics

import (
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestNewMetrics(t *testing.T) {
	client := resty.New()
	m := NewMetrics(client)
	if m == nil {
		t.Errorf("NewMetrics() = %v, want a valid Metrics object", m)
	}
}

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
			err := m.PollMetrics()
			if err != nil {
				t.Errorf("PollMetrics() error = %v, want nil", err)
			}
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

func TestReportMetrics(t *testing.T) {
	setup()
	tests := []struct {
		name           string
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
		client         *resty.Client
		senderText     MetricSender
		senderJSON     MetricSender
		wantErr        bool
	}{
		// Test cases updated to reflect the new error conditions
		{
			name:           "Valid Metrics",
			gaugeMetrics:   map[string]float64{"test": 1.1234567890123456},
			counterMetrics: map[string]int64{"Pool": 1},
			client:         resty.New(),
			senderText:     &TextMetricSender{},
			senderJSON:     &JSONMetricSender{},
			wantErr:        false,
		},
		{
			name:           "Nil GaugeMetrics",
			gaugeMetrics:   nil,
			counterMetrics: map[string]int64{"Pool": 1},
			client:         resty.New(),
			senderText:     &TextMetricSender{},
			senderJSON:     &JSONMetricSender{},
			wantErr:        true,
		},
		{
			name:           "Nil CounterMetrics",
			gaugeMetrics:   map[string]float64{"test": 1.0},
			counterMetrics: nil,
			client:         resty.New(),
			senderText:     &TextMetricSender{},
			senderJSON:     &JSONMetricSender{},
			wantErr:        true,
		},
		{
			name:           "Nil Client",
			gaugeMetrics:   map[string]float64{"test": 1.0},
			counterMetrics: map[string]int64{"Pool": 1},
			client:         nil,
			senderText:     &TextMetricSender{},
			senderJSON:     &JSONMetricSender{},
			wantErr:        true,
		},
		{
			name:           "Nil TextSender",
			gaugeMetrics:   map[string]float64{"test": 1.0},
			counterMetrics: map[string]int64{"Pool": 1},
			client:         resty.New(),
			senderText:     nil,
			senderJSON:     &JSONMetricSender{},
			wantErr:        true,
		},
		{
			name:           "Nil JSONSender",
			gaugeMetrics:   map[string]float64{"test": 1.0},
			counterMetrics: map[string]int64{"Pool": 1},
			client:         resty.New(),
			senderText:     &TextMetricSender{},
			senderJSON:     nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ma := &Metrics{
				GaugeMetrics:   tt.gaugeMetrics,
				CounterMetrics: tt.counterMetrics,
				Client:         tt.client,
				TextSender:     tt.senderText,
				JSONSender:     tt.senderJSON,
			}

			if err := ma.ReportMetrics("http://localhost"); (err != nil) != tt.wantErr {
				t.Errorf("ReportMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
