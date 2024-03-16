package config

// Here we define two global variables to store the metrics we collect.
// We use regular maps, cause they get filled one by one, not concurrently.
var (
	GaugeMetrics   = make(map[string]float64) // Метрики типа gauge
	CounterMetrics = make(map[string]int64)   // Метрики типа counter
)

