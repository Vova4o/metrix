package clientmetrics

import (
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Metric struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MetricSender interface {
	SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error
}

type TextMetricSender struct {
	mu sync.Mutex
}

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type JSONMetricSender struct {
	mu sync.Mutex
}

type MetricsClient interface {
	PollMetrics() error
	ReportMetrics(baseURL string) error
}

type Metrics struct {
	mu             sync.Mutex
	GaugeMetrics   map[string]float64 `json:"gauge"`
	CounterMetrics map[string]int64   `json:"counter"`
	Client         *resty.Client
	PollTicker     *time.Ticker
	ReportTicker   *time.Ticker
	BaseURL        string
	TextSender     MetricSender
	JSONSender     MetricSender
}
