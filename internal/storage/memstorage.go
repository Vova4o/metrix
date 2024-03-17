package storage

import (
	"sync"
)

type MemStorage struct {
	mu             sync.Mutex
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
	Err            error
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
		Err:            nil,
	}
}

func (ms *MemStorage) GetAllGauges() map[string]float64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.GaugeMetrics
}

func (ms *MemStorage) GetAllCounters() map[string]int64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.CounterMetrics
}

func (ms *MemStorage) SetGauge(key string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.GaugeMetrics[key] = value
}

func (ms *MemStorage) GetGauge(key string) (float64, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	value, exists := ms.GaugeMetrics[key]
	return value, exists
}

func (ms *MemStorage) SetCounter(key string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.CounterMetrics[key] += value
}

func (ms *MemStorage) GetCounter(key string) (int64, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	value, exists := ms.CounterMetrics[key]
	return value, exists
}

// func (ms *MemStorage) Delete(key string) {
// 	ms.mu.Lock()
// 	defer ms.mu.Unlock()

// 	delete(ms.GaugeMetrics, key)
// 	delete(ms.CounterMetrics, key)
// }
