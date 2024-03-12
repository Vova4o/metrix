package config

// Metric is a generic metric type that can be used for any type of metric
// It is used to serialize and deserialize metrics to and from JSON
type Metric struct {
	Type  string  `json:"type"`
	Key   string  `json:"key"`
	Value float64 `json:"value"`
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
