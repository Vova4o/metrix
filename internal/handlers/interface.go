package handlers

import (
	"fmt"
	"strconv"

	"Vova4o/metrix/internal/storage"
)

// MetricType is an interface for metric types
type Metricer interface {
	ParseValue(string) (interface{}, error)
	GetValue(storage.Storager, string) (interface{}, bool)
	FormatValue(interface{}) string
	Store(storage.Storager, string, interface{})
	GetAll(storage.Storager) map[string]interface{}
}

type GaugeMetricType struct{}

type CounterMetricType struct{}

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricUpdate struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (g GaugeMetricType) GetAll(storage storage.Storager) map[string]interface{} {
	gauges := storage.GetAllGauges()
	result := make(map[string]interface{}, len(gauges))
	for k, v := range gauges {
		result[k] = v
	}
	return result
}

func (c CounterMetricType) GetAll(storage storage.Storager) map[string]interface{} {
	counters := storage.GetAllCounters()
	result := make(map[string]interface{}, len(counters))
	for k, v := range counters {
		result[k] = v
	}
	return result
}

func (g GaugeMetricType) ParseValue(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}

func (g GaugeMetricType) Store(storage storage.Storager, name string, value interface{}) {
	storage.SetGauge(name, value.(float64))
}

func (c CounterMetricType) ParseValue(value string) (interface{}, error) {
	return strconv.ParseInt(value, 10, 64)
}

func (c CounterMetricType) Store(storage storage.Storager, name string, value interface{}) {
	storage.SetCounter(name, value.(int64))
}

func (g GaugeMetricType) GetValue(storage storage.Storager, name string) (interface{}, bool) {
	return storage.GetGauge(name)
}

func (g GaugeMetricType) FormatValue(value interface{}) string {
	return strconv.FormatFloat(value.(float64), 'f', -1, 64)
}

func (c CounterMetricType) GetValue(storage storage.Storager, name string) (interface{}, bool) {
	return storage.GetCounter(name)
}

func (c CounterMetricType) FormatValue(value interface{}) string {
	return fmt.Sprintf("%d", int(value.(int64)))
}
