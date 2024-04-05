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

type FileStorager interface {
	Storager
	SaveToFile() error
	LoadFromFile() error
	Close() error
}

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

func storeMetricJSON(s Storager, metric MetricsJSON) error {
	switch metric.MType {
	case "gauge":
		if metric.Value != nil {
			s.SetGauge(metric.ID, *metric.Value)
			return nil
		} else {
			return fmt.Errorf("value is required for gauge type")
		}
	case "counter":
		if metric.Delta != nil {
			s.SetCounter(metric.ID, *metric.Delta)
			return nil
		} else {
			return fmt.Errorf("delta is required for counter type")
		}
	default:
		return fmt.Errorf("invalid metric type: %s", metric.MType)
	}
}

func storeMetric(s Storager, metricType, metricName, metricValue string) error {
	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return fmt.Errorf("invalid metric value: %v", err)
		}
		s.SetGauge(metricName, value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid metric value: %v", err)
		}
		s.SetCounter(metricName, value)
	default:
		return fmt.Errorf("invalid metric type: %s", metricType)
	}

	return nil
}

func getJSONValue(s Storager, metrics MetricsJSON) (interface{}, error) {
	var value interface{}
	var ok bool
	switch metrics.MType {
	case "gauge":
		if value, ok = s.GetGauge(metrics.ID); !ok {
			return nil, fmt.Errorf("metric not found")
		}
	case "counter":
		if value, ok = s.GetCounter(metrics.ID); !ok {
			return nil, fmt.Errorf("metric not found")
		}
	default:
		return nil, fmt.Errorf("invalid metric type: %s", metrics.MType)
	}
	return value, nil
}

func getMetricValue(s Storager, metricType, metricName string) (interface{}, error) {
	var value interface{}
	var ok bool
	switch metricType {
	case "gauge":
		if value, ok = s.GetGauge(metricName); !ok {
			return nil, fmt.Errorf("metric not found")
		}
	case "counter":
		if value, ok = s.GetCounter(metricName); !ok {
			return nil, fmt.Errorf("metric not found")
		}
	default:
		return nil, fmt.Errorf("invalid metric type: %s", metricType)
	}
	return value, nil
}
