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
	"time"

	allflags "Vova4o/metrix/internal/flag"
	"github.com/go-resty/resty/v2"
)

type Metrics struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
	Client         *resty.Client
	PollTicker     *time.Ticker
	ReportTicker   *time.Ticker
	BaseURL        string
}
type MetricsClient interface {
	PollMetrics() error
	ReportMetrics(baseURL string) error
}

func NewMetrics(client *resty.Client) *Metrics {
	return &Metrics{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
		Client:         client,
		PollTicker:     time.NewTicker(time.Duration(allflags.GetPollInterval()) * time.Second),
		ReportTicker:   time.NewTicker(time.Duration(allflags.GetReportInterval()) * time.Second),
		BaseURL:        allflags.GetServerAddress(),
	}
}

func (ma *Metrics) PollMetrics() error {
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

	ma.CounterMetrics["PoolCount"]++

	return nil
}

func (ma *Metrics) ReportMetrics(baseURL string) error {
	if ma.GaugeMetrics == nil {
		return errors.New("random value is nil")
	}
	if ma.CounterMetrics == nil {
		return errors.New("counter metrics is nil")
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

	for metricName, metricValue := range ma.GaugeMetrics {
		wg.Add(1)
		go reportMetric("gauge", metricName, strconv.FormatFloat(metricValue, 'f', -1, 64))
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
		log.Println(err)
	}

	return nil
}

func SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error {
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
