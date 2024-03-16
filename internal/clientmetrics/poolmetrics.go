package clientmetrics

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
)

type RestClient interface {
	R() *resty.Request
}

type MetricsAgent struct {
	Metrics *Metrics
	Client  *resty.Client
}

type Metrics struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}
type MetricsClient interface {
	PollMetrics()
	ReportMetrics(baseURL string)
}

func NewMetrics() *Metrics {
	return &Metrics{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}
}

func (m *Metrics) PollMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.GaugeMetrics = map[string]float64{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": float64(memStats.GCCPUFraction),
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}

	m.CounterMetrics["PoolCount"]++
}

func (ma *MetricsAgent) ReportMetrics(baseURL string) error {
	if ma.Metrics == nil {
		return errors.New("metrics is nil")
	}
	if ma.Client == nil {
		return errors.New("client is nil")
	}

	errs := make(chan error)
	var wg sync.WaitGroup

	reportMetric := func(metricType, name, value string) {
		defer wg.Done()
		if err := SendMetric(ma.Client, metricType, name, value, baseURL); err != nil {
			errs <- fmt.Errorf("error sending %s metric %s: %v", metricType, name, err)
		}
	}

	for metricName, metricValue := range ma.Metrics.GaugeMetrics {
		wg.Add(1)
		go reportMetric("gauge", metricName, strconv.FormatFloat(metricValue, 'f', -1, 64))
	}

	for metricName, metricValue := range ma.Metrics.CounterMetrics {
		wg.Add(1)
		go reportMetric("counter", metricName, strconv.FormatInt(metricValue, 10))
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		log.Println(err)
	}

	return nil
}

func SendMetric(client RestClient, metricType, metricName, metricValue, baseURL string) error {
	if client == nil {
        return errors.New("client is nil")
    }
	
	if !strings.HasPrefix(baseURL, "http://") {
		baseURL = "http://" + baseURL
	}
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/%s/%s/%s", baseURL, metricType, metricName, metricValue))

	if err != nil {
		log.Printf("failed to send %s metric %s: %v", metricType, metricName, err)
		return fmt.Errorf("failed to send %s metric %s: %v", metricType, metricName, err)
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
		return fmt.Errorf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
	}

	return nil
}

// package clientmetrics

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"runtime"
// 	"strings"
// 	"sync"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/sirupsen/logrus"
// )

// // HttpClient is an interface for making HTTP requests
// type RestClient interface {
// 	R() *resty.Request
// }

// type MetricsAgent struct {
// 	*MetricsCollector
// 	Client *resty.Client
// }

// type MetricsCollector struct {
// 	GaugeMetrics   map[string]float64
// 	CounterMetrics map[string]float64
// }

// func NewMetricsCollector() *MetricsCollector {
// 	return &MetricsCollector{
// 		GaugeMetrics:   make(map[string]float64),
// 		CounterMetrics: make(map[string]float64),
// 	}
// }

// type MetricsClient interface {
// 	PollMetrics()
// 	ReportMetrics(baseURL string)
// }

// func (mc *MetricsCollector) PollMetrics() {
// 	var memStats runtime.MemStats
// 	runtime.ReadMemStats(&memStats)

// 	mc.GaugeMetrics = map[string]float64{
// 		"Alloc":         float64(memStats.Alloc),
// 		"BuckHashSys":   float64(memStats.BuckHashSys),
// 		"Frees":         float64(memStats.Frees),
// 		"GCCPUFraction": float64(memStats.GCCPUFraction),
// 		"GCSys":         float64(memStats.GCSys),
// 		"HeapAlloc":     float64(memStats.HeapAlloc),
// 		"HeapIdle":      float64(memStats.HeapIdle),
// 		"HeapInuse":     float64(memStats.HeapInuse),
// 		"HeapObjects":   float64(memStats.HeapObjects),
// 		"HeapReleased":  float64(memStats.HeapReleased),
// 		"HeapSys":       float64(memStats.HeapSys),
// 		"LastGC":        float64(memStats.LastGC),
// 		"Lookups":       float64(memStats.Lookups),
// 		"MCacheInuse":   float64(memStats.MCacheInuse),
// 		"MCacheSys":     float64(memStats.MCacheSys),
// 		"MSpanInuse":    float64(memStats.MSpanInuse),
// 		"MSpanSys":      float64(memStats.MSpanSys),
// 		"Mallocs":       float64(memStats.Mallocs),
// 		"NextGC":        float64(memStats.NextGC),
// 		"NumForcedGC":   float64(memStats.NumForcedGC),
// 		"NumGC":         float64(memStats.NumGC),
// 		"OtherSys":      float64(memStats.OtherSys),
// 		"PauseTotalNs":  float64(memStats.PauseTotalNs),
// 		"StackInuse":    float64(memStats.StackInuse),
// 		"StackSys":      float64(memStats.StackSys),
// 		"Sys":           float64(memStats.Sys),
// 		"TotalAlloc":    float64(memStats.TotalAlloc),
// 		"RandomValue":   rand.Float64(),
// 	}

// 	mc.CounterMetrics["PoolCount"]++
// }

// func (mc *MetricsAgent) ReportMetrics(baseURL string) {
// 	// Add a middleware logger
// 	mc.Client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
// 		logrus.WithFields(logrus.Fields{
// 			"url": request.URL,
// 		}).Info("Sending request")

// 		return nil
// 	})

// 	mc.Client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
// 		logrus.WithFields(logrus.Fields{
// 			"status": response.StatusCode(),
// 			"body":   response.String(),
// 		}).Info("Received response")

// 		return nil
// 	})

// 	errs := make(chan error)
// 	var wg sync.WaitGroup

// 	for metricName, metricValue := range mc.GaugeMetrics {
// 		wg.Add(1)
// 		go func(name string, value float64) {
// 			defer wg.Done()
// 			err := SendMetric(mc.Client, "gauge", name, value, baseURL)
// 			if err != nil {
// 				log.Printf("error sending gauge metric %s: %v", name, err)
// 				errs <- fmt.Errorf("error sending gauge metric %s: %v", name, err)
// 			}
// 		}(metricName, metricValue)
// 	}

// 	for metricName, metricValue := range mc.CounterMetrics {
// 		wg.Add(1)
// 		go func(name string, value float64) {
// 			defer wg.Done()
// 			err := SendMetric(mc.Client, "counter", name, float64(value), baseURL)
// 			if err != nil {
// 				log.Printf("error sending counter metric %s: %v", name, err)
// 				errs <- fmt.Errorf("error sending counter metric %s: %v", name, err)
// 			}
// 		}(metricName, metricValue)
// 	}

// 	// Close the errs channel after all goroutines have finished
// 	go func() {
// 		wg.Wait()
// 		close(errs)
// 	}()

// 	// Print all errors
// 	for err := range errs {
// 		log.Println(err)
// 		fmt.Println(err)
// 	}
// }

// func SendMetric(client RestClient, metricType, metricName string, metricValue float64, baseURL string) error {
// 	if !strings.HasPrefix(baseURL, "http://") {
// 		baseURL = "http://" + baseURL
// 	}
// 	// fmt.Println("Sending metric", metricType, metricName, metricValue, baseURL)
// 	resp, err := client.R().
// 		SetHeader("Content-Type", "text/plain").
// 		Post(fmt.Sprintf("%s/update/%s/%s/%.2f", baseURL, metricType, metricName, metricValue))

// 	if err != nil {
// 		log.Printf("failed to send %s metric %s: %v", metricType, metricName, err)
// 		return fmt.Errorf("failed to send %s metric %s: %v", metricType, metricName, err)
// 	}

// 	if resp.StatusCode() != http.StatusOK {
// 		log.Printf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
// 		return fmt.Errorf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
// 	}

// 	return nil
// }

// // package clientmetrics

// // import (
// // 	"math/rand"
// // 	"runtime"

// // 	"Vova4o/metrix/internal/config"
// // )

// // // pollMetrics collects the metrics
// // func PollMetrics() {
// // 	// Собираем метрики
// // 	var memStats runtime.MemStats
// // 	runtime.ReadMemStats(&memStats)

// // 	// Обновляем метрики типа gauge
// // 	config.GaugeMetrics = map[string]float64{
// // 		"Alloc":         float64(memStats.Alloc),
// // 		"BuckHashSys":   float64(memStats.BuckHashSys),
// // 		"Frees":         float64(memStats.Frees),
// // 		"GCCPUFraction": float64(memStats.GCCPUFraction),
// // 		"GCSys":         float64(memStats.GCSys),
// // 		"HeapAlloc":     float64(memStats.HeapAlloc),
// // 		"HeapIdle":      float64(memStats.HeapIdle),
// // 		"HeapInuse":     float64(memStats.HeapInuse),
// // 		"HeapObjects":   float64(memStats.HeapObjects),
// // 		"HeapReleased":  float64(memStats.HeapReleased),
// // 		"HeapSys":       float64(memStats.HeapSys),
// // 		"LastGC":        float64(memStats.LastGC),
// // 		"Lookups":       float64(memStats.Lookups),
// // 		"MCacheInuse":   float64(memStats.MCacheInuse),
// // 		"MCacheSys":     float64(memStats.MCacheSys),
// // 		"MSpanInuse":    float64(memStats.MSpanInuse),
// // 		"MSpanSys":      float64(memStats.MSpanSys),
// // 		"Mallocs":       float64(memStats.Mallocs),
// // 		"NextGC":        float64(memStats.NextGC),
// // 		"NumForcedGC":   float64(memStats.NumForcedGC),
// // 		"NumGC":         float64(memStats.NumGC),
// // 		"OtherSys":      float64(memStats.OtherSys),
// // 		"PauseTotalNs":  float64(memStats.PauseTotalNs),
// // 		"StackInuse":    float64(memStats.StackInuse),
// // 		"StackSys":      float64(memStats.StackSys),
// // 		"Sys":           float64(memStats.Sys),
// // 		"TotalAlloc":    float64(memStats.TotalAlloc),
// // 		"RandomValue":   rand.Float64(), // Some random value
// // 	}

// // 	// Обновляем метрики типа counter
// // 	config.CounterMetrics["PoolCount"]++

// // }
