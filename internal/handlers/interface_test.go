package handlers

import (
	"log"
	"testing"
)

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

func (m *mockStorager) Store(key string, value interface{}) {
	switch v := value.(type) {
	case float64:
		m.SetGauge(key, v)
	case int64:
		m.SetCounter(key, v)
	default:
		log.Printf("Invalid value type: %T", v)
	}
}

func TestGaugeMetricType_GetAll(t *testing.T) {
	mock := &mockStorager{
		gauges: map[string]float64{
			"gauge1": 10.5,
			"gauge2": 20.5,
		},
		counters: map[string]int64{},
	}

	g := GaugeMetricType{}
	all := g.GetAll(mock)

	if len(all) != 2 {
		t.Errorf("expected %v, got %v", 2, len(all))
	}
	if all["gauge1"] != 10.5 {
		t.Errorf("expected %v, got %v", 10.5, all["gauge1"])
	}
	if all["gauge2"] != 20.5 {
		t.Errorf("expected %v, got %v", 20.5, all["gauge2"])
	}
}

func TestCounterMetricType_GetAll(t *testing.T) {
	mock := &mockStorager{
		gauges: map[string]float64{},
		counters: map[string]int64{
			"counter1": 100,
			"counter2": 200,
		},
	}

	c := CounterMetricType{}
	all := c.GetAll(mock)

	if len(all) != 2 {
		t.Errorf("expected %v, got %v", 2, len(all))
	}
	counter1, ok := all["counter1"].(int64) // assert that all["counter1"] is an int64
	if !ok {
		t.Errorf("all[\"counter1\"] is not an int64: %v", all["counter1"])
	} else if counter1 != 100 {
		t.Errorf("expected %v, got %v", 100, counter1)
	}
	counter2, ok := all["counter2"].(int64) // assert that all["counter2"] is an int64
	if !ok {
		t.Errorf("all[\"counter2\"] is not an int64: %v", all["counter2"])
	} else if counter2 != 200 {
		t.Errorf("expected %v, got %v", 200, counter2)
	}
}

func TestGaugeMetricType_ParseValue(t *testing.T) {
	g := GaugeMetricType{}
	value, err := g.ParseValue("10.5")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if value != 10.5 {
		t.Errorf("expected %v, got %v", 10.5, value)
	}
}

func TestGaugeMetricType_Store(t *testing.T) {
	mock := &mockStorager{
		gauges:   map[string]float64{},
		counters: map[string]int64{},
	}

	g := GaugeMetricType{}
	g.Store(mock, "gauge1", 10.5)

	if len(mock.gauges) != 1 {
		t.Errorf("expected %v, got %v", 1, len(mock.gauges))
	}
	if mock.gauges["gauge1"] != 10.5 {
		t.Errorf("expected %v, got %v", 10.5, mock.gauges["gauge1"])
	}
}

func TestCounterMetricType_ParseValue(t *testing.T) {
	c := CounterMetricType{}
	value, err := c.ParseValue("100")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	intValue, ok := value.(int64)
	if !ok {
		t.Errorf("value is not an int: %v", value)
	} else if intValue != 100 {
		t.Errorf("expected %v, got %v", 100, intValue)
	}
}

func TestCounterMetricType_Store(t *testing.T) {
	mock := &mockStorager{
		gauges:   map[string]float64{},
		counters: map[string]int64{},
	}

	c := CounterMetricType{}
	c.Store(mock, "counter1", int64(100)) // convert int to int64

	if len(mock.counters) != 1 {
		t.Errorf("expected %v, got %v", 1, len(mock.counters))
	}
	if mock.counters["counter1"] != int64(100) { // compare with int64(100)
		t.Errorf("expected %v, got %v", int64(100), mock.counters["counter1"])
	}
}

func TestGaugeMetricType_GetValue(t *testing.T) {
	mock := &mockStorager{
		gauges: map[string]float64{
			"gauge1": 10.5,
		},
		counters: map[string]int64{},
	}

	g := GaugeMetricType{}
	value, ok := g.GetValue(mock, "gauge1")

	if !ok {
		t.Errorf("expected %v, got %v", true, ok)
	}
	if value != 10.5 {
		t.Errorf("expected %v, got %v", 10.5, value)
	}
}

func TestCounterMetricType_GetValue(t *testing.T) {
	mock := &mockStorager{
		gauges: map[string]float64{},
		counters: map[string]int64{
			"counter1": 100,
		},
	}

	c := CounterMetricType{}
	value, ok := c.GetValue(mock, "counter1")

	if !ok {
		t.Errorf("expected %v, got %v", true, ok)
	}
	intValue, ok := value.(int64) // assert that value is an int64
	if !ok {
		t.Errorf("value is not an int64: %v", value)
	} else if intValue != 100 {
		t.Errorf("expected %v, got %v", 100, intValue)
	}
}

func TestGaugeMetricType_FormatValue(t *testing.T) {
	g := GaugeMetricType{}
	formatted := g.FormatValue(10.5)

	if formatted != "10.5" {
		t.Errorf("expected %v, got %v", "10.5", formatted)
	}
}

func TestCounterMetricType_FormatValue(t *testing.T) {
	c := CounterMetricType{}
	formatted := c.FormatValue(int64(100))

	if formatted != "100" {
		t.Errorf("expected %v, got %v", "100", formatted)
	}
}
