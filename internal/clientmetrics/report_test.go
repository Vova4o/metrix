package clientmetrics

// import (
// 	"testing"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/stretchr/testify/assert"
// )

// type MockRestClient struct {
// 	*resty.Client
// }

// func (m *MockRestClient) R() *resty.Request {
// 	return m.Client.R() // Use the R() method of resty.Client
// }

// // Mock SendMetric function
// func MockSendMetric(t *testing.T, client *resty.Client, metricType, name string, value float64, baseURL string) error {
//     // Check the parameters
//     assert.Equal(t, "gauge", metricType)
//     assert.Equal(t, "gaugeTest", name)
//     assert.Equal(t, 10.0000, value)
//     assert.Equal(t, "http://localhost:8080", baseURL)

//     // Return nil to indicate no error
//     return nil
// }

// func TestReportMetrics(t *testing.T) {
//     // Set up the mock client
//     client := &MockRestClient{Client: resty.New()}

//     // Set up the metrics
//     mc := &MetricsCollector{
//         GaugeMetrics: map[string]float64{
//             "gaugeTest": 10.0000,
//         },
//         CounterMetrics: map[string]float64{
//             "counterTest": 20,
//         },
//     }

//     // Use mc in your test
//     err := mc.ReportMetrics(client, "http://localhost:8080")
//     assert.Nil(t, err)
// }
