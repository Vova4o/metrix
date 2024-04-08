package storage

import (
	"sync"

	"Vova4o/metrix/internal/handlers"
)

// MemStorage is a simple in-memory storage
// that implements the StorageInterface
// It uses a mutex to synchronize access to the maps
// GaugeMetrics and CounterMetrics
type MemStorage struct {
	mu             sync.Mutex
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
	Err            error
}

// NewMemStorage creates a new MemStorage
// and returns a pointer to it
// GaugeMetrics and CounterMetrics are initialized as empty maps
func NewMemory() handlers.Storager {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
		Err:            nil,
	}
}

// GetAllGauges returns a map of all gauge metrics
func (ms *MemStorage) GetAllGauges() map[string]float64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.GaugeMetrics
}

// GetAllCounters returns a map of all counter metrics
func (ms *MemStorage) GetAllCounters() map[string]int64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.CounterMetrics
}

// SetGauge sets the value of a gauge metric
func (ms *MemStorage) SetGauge(key string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.GaugeMetrics[key] = value
}

// GetGauge returns the value of a gauge metric
func (ms *MemStorage) GetGauge(key string) (float64, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	value, exists := ms.GaugeMetrics[key]
	return value, exists
}

// SetCounter sets the value of a counter metric
func (ms *MemStorage) SetCounter(key string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.CounterMetrics[key] += value
}

// GetCounter returns the value of a counter metric
func (ms *MemStorage) GetCounter(key string) (int64, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	value, exists := ms.CounterMetrics[key]
	return value, exists
}

func (ms *MemStorage) GetAllMetrics() map[string]interface{} {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return map[string]interface{}{
		"Gauge":   ms.GaugeMetrics,
		"Counter": ms.CounterMetrics,
	}
}

// func (ms *MemStorage) Delete(key string) {
// 	ms.mu.Lock()
// 	defer ms.mu.Unlock()

// 	delete(ms.GaugeMetrics, key)
// 	delete(ms.CounterMetrics, key)
// }
