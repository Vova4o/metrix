package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// HttpClient is an interface for making HTTP requests
type RestClient interface {
	R() *resty.Request
}

var baseURL string

// Here we define two global variables to store the metrics we collect.
// We use regular maps, cause they get filled one by one, not concurrently.
var (
	gaugeMetrics   = make(map[string]float64) // Метрики типа gauge
	counterMetrics = make(map[string]int64)   // Метрики типа counter
)

func main() {
	// Open a file for logging
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Set the output destination of the standard logger
	log.SetOutput(logFile)

	// Parse the flags
	parseFlags()

	//Start the cicle of collecting and sending metrics
	pollTicker := time.NewTicker(*PollInterval)
	reportTicker := time.NewTicker(*ReportInterval)
	baseURL = *ServerAddress

	// Start the main loop
	for {
		// Wait for the next tick
		select {
		// When the pollTicker ticks, we collect the metrics
		case <-pollTicker.C:
			pollMetrics()
			// When the reportTicker ticks, we send the metrics
		case <-reportTicker.C:
			reportMetrics(baseURL)
		}
	}
}

// pollMetrics collects the metrics
func pollMetrics() {
	// Собираем метрики
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Обновляем метрики типа gauge
	gaugeMetrics = map[string]float64{
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
		"RandomValue":   rand.Float64(), // Some random value
	}

	// Обновляем метрики типа counter
	counterMetrics["PoolCount"]++

}

// reportMetrics sends the metrics to the server
func reportMetrics(baseURL string) {
	client := resty.New()
	errs := make(chan error)
	var wg sync.WaitGroup

	for metricName, metricValue := range gaugeMetrics {
		wg.Add(1)
		go func(name string, value float64) {
			defer wg.Done()
			err := sendMetric(client, "gauge", name, value, baseURL)
			if err != nil {
				log.Printf("error sending gauge metric %s: %v", name, err)
				errs <- fmt.Errorf("error sending gauge metric %s: %v", name, err)
			}
		}(metricName, metricValue)
	}

	for metricName, metricValue := range counterMetrics {
		wg.Add(1)
		go func(name string, value float64) {
			defer wg.Done()
			err := sendMetric(client, "counter", name, value, baseURL)
			if err != nil {
				log.Printf("error sending counter metric %s: %v", name, err)
				errs <- fmt.Errorf("error sending counter metric %s: %v", name, err)
			}
		}(metricName, float64(metricValue))
	}

	// Close the errs channel after all goroutines have finished
	go func() {
		wg.Wait()
		close(errs)
	}()

	// Print all errors
	for err := range errs {
		log.Println(err)
		fmt.Println(err)
	}
}

func sendMetric(client RestClient, metricType, metricName string, metricValue float64, baseURL string) error {
	if !strings.HasPrefix(baseURL, "http://") {
		baseURL = "http://" + baseURL
	}
	// fmt.Println("Sending metric", metricType, metricName, metricValue, baseURL)
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/%s/%s/%.2f", baseURL, metricType, metricName, metricValue))

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
