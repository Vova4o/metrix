package clientmetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

type JSONMetricSender struct{}

func (j *JSONMetricSender) SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error {
	metric := Metric{
		Type:  metricType,
		Name:  metricName,
		Value: metricValue,
	}

	if client == nil {
		return errors.New("client is nil")
	}

	if !strings.HasPrefix(baseURL, "http://") {
		baseURL = "http://" + baseURL
	}

	jsonDate, err := json.Marshal(metric)
	if err != nil {
		log.Printf("failed to marshal metric: %v", err)
		return fmt.Errorf("failed to marshal metric: %v", err)
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonDate).
		Post(fmt.Sprintf("%s/update/json/", baseURL))

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
