package handlers

import (
	"fmt"
	"strconv"
)

type Storager interface {
	SetGauge(key string, value float64)
	GetGauge(key string) (float64, bool)
	SetCounter(key string, value int64)
	GetCounter(key string) (int64, bool)
	// Delete(key string)
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
	GetAllMetrics() map[string]interface{}
}

// MetricType is an interface for metric types
type Metricer interface {
	ParseValue(string) (interface{}, error)
	GetValue(Storager, string) (interface{}, bool)
	FormatValue(interface{}) string
	Store(Storager, string, interface{})
	GetAll(Storager) map[string]interface{}
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
	Type  string
	Name  string
	Value string
}

func (g GaugeMetricType) GetAll(s Storager) map[string]interface{} {
	gauges := s.GetAllGauges()
	result := make(map[string]interface{}, len(gauges))
	for k, v := range gauges {
		result[k] = v
	}
	return result
}

func (c CounterMetricType) GetAll(s Storager) map[string]interface{} {
	counters := s.GetAllCounters()
	result := make(map[string]interface{}, len(counters))
	for k, v := range counters {
		result[k] = v
	}
	return result
}

func (g GaugeMetricType) ParseValue(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}

func (g GaugeMetricType) Store(s Storager, name string, value interface{}) {
	s.SetGauge(name, value.(float64))
}

func (c CounterMetricType) ParseValue(value string) (interface{}, error) {
	return strconv.ParseInt(value, 10, 64)
}

func (c CounterMetricType) Store(s Storager, name string, value interface{}) {
	s.SetCounter(name, value.(int64))
}

func (g GaugeMetricType) GetValue(s Storager, name string) (interface{}, bool) {
	return s.GetGauge(name)
}

func (g GaugeMetricType) FormatValue(value interface{}) string {
	return strconv.FormatFloat(value.(float64), 'f', -1, 64)
}

func (c CounterMetricType) GetValue(s Storager, name string) (interface{}, bool) {
	return s.GetCounter(name)
}

func (c CounterMetricType) FormatValue(value interface{}) string {
	return fmt.Sprintf("%d", int(value.(int64)))
}
