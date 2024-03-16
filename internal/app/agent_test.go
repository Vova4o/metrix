package app

// import (
// 	"Vova4o/metrix/internal/clientmetrics"
// 	"testing"

// 	"github.com/go-resty/resty/v2"
// )

// func TestHandleTick(t *testing.T) {
// 	// Create a mock MetricsAgent
// 	ma := &clientmetrics.MetricsAgent{
// 		Metrics: &clientmetrics.Metrics{},
// 		Client:  &resty.Client{},
// 	}

// 	tests := []struct {
// 		name     string
// 		tickType string
// 	}{
// 		{
// 			name:     "Test Case 1 - Poll tick",
// 			tickType: "poll",
// 		},
// 		{
// 			name:     "Test Case 2 - Report tick",
// 			tickType: "report",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			HandleTick(ma, "http://example.com", tt.tickType)
// 			// Add assertions here to check the behavior of HandleTick
// 		})
// 	}
// }
