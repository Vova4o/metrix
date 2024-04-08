package handlers

type mockStorager struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *mockStorager) SetGauge(name string, value float64) {
	if m.gauges == nil {
		m.gauges = make(map[string]float64)
	}
	m.gauges[name] = value
}

func (m *mockStorager) GetGauge(key string) (float64, bool) {
	value, ok := m.gauges[key]
	return value, ok
}

func (m *mockStorager) SetCounter(name string, value int64) {
	if m.counters == nil {
		m.counters = make(map[string]int64)
	}
	m.counters[name] = value
}

func (m *mockStorager) GetCounter(key string) (int64, bool) {
	value, ok := m.counters[key]
	return value, ok
}

func (m *mockStorager) GetAllGauges() map[string]float64 {
	return m.gauges
}

func (m *mockStorager) GetAllCounters() map[string]int64 {
	return m.counters
}

func (m *mockStorager) GetAllMetrics() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m.gauges {
		result[k] = v
	}
	for k, v := range m.counters {
		result[k] = v
	}
	return result
}

func (m *mockStorager) GetValue(metricType, metricName string) (float64, bool) {
	if metricType == "gauge" {
		value, ok := m.gauges[metricName]
		return value, ok
	} else if metricType == "counter" {
		value, ok := m.counters[metricName]
		return float64(value), ok
	}
	return 0, false
}

// func (m *mockStorager) Store(key string, value interface{}) {
// 	switch v := value.(type) {
// 	case float64:
// 		m.SetGauge(key, v)
// 	case int64:
// 		m.SetCounter(key, v)
// 	default:
// 		log.Printf("Invalid value type: %T", v)
// 	}
// }
