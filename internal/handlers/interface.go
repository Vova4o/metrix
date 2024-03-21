package handlers

import "Vova4o/metrix/internal/storage"

// MetricType is an interface for metric types
type MetricType interface {
	ParseValue(string) (interface{}, error)
	GetValue(storage.StorageInterface, string) (interface{}, bool)
	FormatValue(interface{}) string
	Store(storage.StorageInterface, string, interface{})
	GetAll(storage.StorageInterface) map[string]interface{}
}

type GaugeMetricType struct{}

type CounterMetricType struct{}

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
