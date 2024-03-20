package clientmetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Metric struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MetricSender interface {
	SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error
}

type TextMetricSender struct{}

func (t *TextMetricSender) SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error {
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

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type JSONMetricSender struct{}

func (j *JSONMetricSender) SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error {
	var delta *int64
	var value *float64

	if metricType == "counter" {
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse counter value: %v", err)
		}
		delta = &val
	} else if metricType == "gauge" {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return fmt.Errorf("failed to parse gauge value: %v", err)
		}
		value = &val
	} else {
		return fmt.Errorf("invalid metric type: %s", metricType)
	}

	metric := MetricsJSON{
		ID:    metricName,
		MType: metricType,
		Delta: delta,
		Value: value,
	}

	if client == nil {
		return errors.New("client is nil")
	}

	if !strings.HasPrefix(baseURL, "http://") {
		baseURL = "http://" + baseURL
	}

	jsonDate, err := json.MarshalIndent(metric, "", "  ")
	if err != nil {
		log.Printf("failed to marshal metric: %v", err)
		return fmt.Errorf("failed to marshal metric: %v", err)
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonDate).
		Post(fmt.Sprintf("%s/update/", baseURL))
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
