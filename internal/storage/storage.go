package storage

type StorageInterface interface {
	SetGauge(key string, value float64)
	GetGauge(key string) (float64, bool)
	SetCounter(key string, value int64)
	GetCounter(key string) (int64, bool)
	// Delete(key string)
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
}
