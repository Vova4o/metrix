package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

// some constants to run the agent
const (
	serverURL      = "http://localhost:8080/update" // URL сервера
	pollInterval   = 2 * time.Second                // Интервал между сбором метрик
	reportInterval = 10 * time.Second               // Интервал между отправкой метрик
)

// Here we define two global variables to store the metrics we collect.
// We use regular maps, cause they get filled one by one, not concurrently.
var (
	gaugeMetrics   = make(map[string]float64) // Метрики типа gauge
	counterMetrics = make(map[string]int64)   // Метрики типа counter
)

func main() {
	//Start the cicle of collecting and sending metrics
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	// Start the main loop
	for {
		// Wait for the next tick
		select {
		// When the pollTicker ticks, we collect the metrics
		case <-pollTicker.C:
			pollMetrics()
			// When the reportTicker ticks, we send the metrics
		case <-reportTicker.C:
			reportMetrics()
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
func reportMetrics() {
	// Send the metrics to the server in parallel to make it faster
	// type of gauge
	for metricName, metricValue := range gaugeMetrics {
		go sendMetric("gauge", metricName, metricValue)
	}

	// Send the metrics to the server in parallel to make it faster
	// type of counter, well its only one for now but we can add more
	for metricName, metricValue := range counterMetrics {
		go sendMetric("counter", metricName, float64(metricValue))
	}
}

// sendMetric sends a metric to the server
func sendMetric(metricType, metricName string, metricValue float64) {
	//Prepare the URL for the request to the server with the metric data
	url := fmt.Sprintf("%s/%s/%s/%s", serverURL, metricType, metricName, strconv.FormatFloat(metricValue, 'f', -1, 64))

	//Send the metric to the server
	//Post issues a POST to the specified URL, as a text/plain, with a nil body.
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Printf("Failed to send %s metric %s: %v\n", metricType, metricName, err)
		return
	}
	//Close the response body
	defer resp.Body.Close()

	//Check if the server returned a non-OK status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Server returned non-OK status for %s metric %s: %v\n", metricType, metricName, resp.Status)
	}
}
