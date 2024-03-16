package clientmetrics

import "github.com/go-resty/resty/v2"

// HttpClient is an interface for making HTTP requests
type RestClient interface {
	R() *resty.Request
}

type MetricsAgent struct {
	*MetricsCollector
	Client *resty.Client
}

type MetricsCollector struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int),
	}
}

// type MetricsClient interface {
// 	PollMetrics()
// 	ReportMetrics(baseURL string)
// }

// type MetricsCollector struct {
// 	GaugeMetrics   map[string]float64
// 	CounterMetrics map[string]int
// 	Client         RestClient
// }

// func NewMetricsCollector(client RestClient) *MetricsCollector {
// 	return &MetricsCollector{
// 		GaugeMetrics:   make(map[string]float64),
// 		CounterMetrics: make(map[string]int),
// 		Client:         client,
// 	}
// }
