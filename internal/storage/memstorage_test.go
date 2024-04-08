package storage

import "testing"

func TestMemStorage_GetAllGauges(t *testing.T) {
	ms := NewMemory()

	ms.SetGauge("gauge1", 1.23)
	ms.SetGauge("gauge2", 4.56)
	ms.SetCounter("counter1", 10)


	// Call the GetAllGauges method
	gauges := ms.GetAllGauges()

	// Check the returned values
	if len(gauges) != 2 {
		t.Errorf("expected %v, got %v", 2, len(gauges))
	}
	if gauges["gauge1"] != 1.23 {
		t.Errorf("expected %v, got %v", 1.23, gauges["gauge1"])
	}
	if gauges["gauge2"] != 4.56 {
		t.Errorf("expected %v, got %v", 4.56, gauges["gauge2"])
	}
}

func TestMemStorage_GetAllCounters(t *testing.T) {
	ms := NewMemory()

	// Set counter metrics
	ms.SetGauge("gauge1", 1.23)
	ms.SetCounter("counter1", 10)
	ms.SetCounter("counter1", 10)

	// Call the GetAllCounters method
	counters := ms.GetAllCounters()

	
	if counters["counter1"] != 20 {
		t.Errorf("expected %v, got %v", 20, counters["counter2"])
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	ms := NewMemory()

	// Call the SetGauge method
	ms.SetGauge("gauge1", 1.23)

	// Check the value of the gauge metric
	value, exists := ms.GetGauge("gauge1")
	if value != 1.23 {
		t.Errorf("expected %v, got %v", 1.23, value)
	}
	if !exists {
		t.Errorf("expected %v, got %v", true, exists)
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	ms := NewMemory()

	ms.SetGauge("gauge1", 1.23)
	ms.SetCounter("counter1", 10)

	// Call the GetGauge method
	value, exists := ms.GetGauge("gauge1")

	// Check the returned values
	if value != 1.23 {
		t.Errorf("expected %v, got %v", 1.23, value)
	}
	if !exists {
		t.Errorf("expected %v, got %v", true, exists)
	}

	// Call the GetGauge method with a non-existent key
	value, exists = ms.GetGauge("nonexistent")

	// Check the returned values
	if value != 0 {
		t.Errorf("expected %v, got %v", 0, value)
	}
	if exists {
		t.Errorf("expected %v, got %v", false, exists)
	}
}

func TestMemStorage_SetCounter(t *testing.T) {
    testCases := []struct {
        name     string
        counter  string
        value    int64
        expected int64
    }{
        {"Test1", "counter", 10, 10},
		{"Test2", "counter", 20, 30},
        // Add more test cases here...
    }
	
	ms := NewMemory()

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {

            // Call the SetCounter method
            ms.SetCounter(tc.counter, tc.value)

            // Check the value of the counter metric
            value, exists := ms.GetCounter(tc.counter)
            if value != tc.expected {
                t.Errorf("expected %v, got %v", tc.expected, value)
            }
            if !exists {
                t.Errorf("expected %v, got %v", true, exists)
            }
        })
    }
}

func TestMemStorage_GetCounter(t *testing.T) {
	ms := NewMemory()

	ms.SetGauge("gauge1", 1.23)
	ms.SetCounter("counter1", 10)

	// Call the GetCounter method
	value, exists := ms.GetCounter("counter1")

	// Check the returned values
	if value != 10 {
		t.Errorf("expected %v, got %v", 10, value)
	}
	if !exists {
		t.Errorf("expected %v, got %v", true, exists)
	}

	// Call the GetCounter method with a non-existent key
	value, exists = ms.GetCounter("nonexistent")

	// Check the returned values
	if value != 0 {
		t.Errorf("expected %v, got %v", 0, value)
	}
	if exists {
		t.Errorf("expected %v, got %v", false, exists)
	}
}

func TestMemStorage_GetAllMetrics(t *testing.T) {
	ms := NewMemory()

	ms.SetGauge("gauge1", 1.23)
	ms.SetCounter("counter1", 10)

	// Call the GetAllMetrics method
	metrics := ms.GetAllMetrics()

	// Check the returned values
	if len(metrics) != 2 {
		t.Errorf("expected %v, got %v", 2, len(metrics))
	}
	if gauges, ok := metrics["Gauge"].(map[string]float64); ok {
		if len(gauges) != 1 {
			t.Errorf("expected %v, got %v", 1, len(gauges))
		}
		if gauges["gauge1"] != 1.23 {
			t.Errorf("expected %v, got %v", 1.23, gauges["gauge1"])
		}
	} else {
		t.Errorf("expected %v, got %v", "map[string]float64", metrics["Gauge"])
	}
	if counters, ok := metrics["Counter"].(map[string]int64); ok {
		if len(counters) != 1 {
			t.Errorf("expected %v, got %v", 1, len(counters))
		}
		if counters["counter1"] != 10 {
			t.Errorf("expected %v, got %v", 10, counters["counter1"])
		}
	} else {
		t.Errorf("expected %v, got %v", "map[string]int64", metrics["Counter"])
	}
}
