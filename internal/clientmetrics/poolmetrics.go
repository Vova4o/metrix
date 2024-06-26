package clientmetrics

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"Vova4o/metrix/internal/agentflags"
	"Vova4o/metrix/internal/logger"

	"github.com/go-resty/resty/v2"
)

func NewMetrics(client *resty.Client) *Metrics {
	return &Metrics{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
		Client:         client,
		PollTicker:     time.NewTicker(time.Duration(agentflags.GetPollInterval()) * time.Second),
		ReportTicker:   time.NewTicker(time.Duration(agentflags.GetReportInterval()) * time.Second),
		BaseURL:        agentflags.GetServerAddress(),
		TextSender:     &TextMetricSender{},
		JSONSender:     &JSONMetricSender{},
	}
}

func (ma *Metrics) PollMetrics() error {
	ma.mu.Lock()
	defer ma.mu.Unlock()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	ma.GaugeMetrics = map[string]float64{
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

	ma.CounterMetrics["PollCount"]++

	return nil
}

func (ma *Metrics) ReportMetrics(baseURL string) error {
	ma.mu.Lock()
	defer ma.mu.Unlock()
	
	if ma.GaugeMetrics == nil {
		return errors.New("random value is nil")
	}
	if ma.CounterMetrics == nil {
		return errors.New("counter metrics is nil")
	}
	if ma.Client == nil {
		return errors.New("client is nil")
	}
	if ma.TextSender == nil {
		return errors.New("text sender is nil")
	}
	if ma.JSONSender == nil {
		return errors.New("json sender is nil")
	}

	errs := make(chan error)
	var wg sync.WaitGroup

	reportMetric := func(metricType, name, value string) {
		defer wg.Done()
		if err := ma.TextSender.SendMetric(ma.Client, metricType, name, value, baseURL); err != nil {
			logger.Log.Errorf("error sending %s metric %s: %v", metricType, name, err)
		}
		if err := ma.JSONSender.SendMetric(ma.Client, metricType, name, value, baseURL); err != nil {
			logger.Log.Errorf("error sending %s metric %s: %v", metricType, name, err)
		}
	}

	for metricName, metricValue := range ma.GaugeMetrics {
		wg.Add(1)
		go reportMetric("gauge", metricName, fmt.Sprintf("%g", metricValue))
	}

	for metricName, metricValue := range ma.CounterMetrics {
		wg.Add(1)
		go reportMetric("counter", metricName, strconv.FormatInt(metricValue, 10))
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		logger.Log.Println(err)
	}

	return nil
}
