package storage

import "sync"

// MemStorage is a simple in-memory storage for metrics
// It uses two sync.Map to store gauge and counter metrics
// I had to change the type of the map to avoid concurrent map writes
type MemStorage struct {
	GaugeMetrics   sync.Map
	CounterMetrics sync.Map
}

// StorageInterface is an interface for storage backends
type StorageInterface interface {
	//SetGauge sets the value of a gauge
	SetGauge(key string, value float64)
	//GetGauge returns the value of a gauge
	GetGauge(key string) (float64, bool)
	//SetCounter sets the value of a counter
	SetCounter(key string, value float64)
	//GetCounter returns the value of a counter
	GetCounter(key string) (float64, bool)
	//Delete removes a metric from the storage
	Delete(key string)
}

// SetGauge sets the value of a gauge
func (m *MemStorage) SetGauge(key string, value float64) {
	m.GaugeMetrics.Store(key, value)
}

// GetGauge returns the value of a gauge
func (m *MemStorage) GetGauge(key string) (float64, bool) {
	value, exists := m.GaugeMetrics.Load(key)
	if exists {
		return value.(float64), exists
	}
	return 0, exists
}

// SetCounter increments the value of a counter
func (m *MemStorage) SetCounter(key string, value float64) {
	actual, loaded := m.CounterMetrics.LoadOrStore(key, value)
	if loaded {
		newValue := actual.(float64) + value
		m.CounterMetrics.Store(key, newValue)
	}
}

// GetCounter returns the value of a counter
func (m *MemStorage) GetCounter(key string) (float64, bool) {
	value, exists := m.CounterMetrics.Load(key)
	if exists {
		return value.(float64), exists
	}
	return 0, exists
}

// Delete removes a metric from the storage
func (m *MemStorage) Delete(key string) {
	m.GaugeMetrics.Delete(key)
	m.CounterMetrics.Delete(key)
}
