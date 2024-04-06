package clientmetrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"Vova4o/metrix/internal/logger"

	"github.com/go-resty/resty/v2"
)

func (t *TextMetricSender) SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if client == nil {
		return errors.New("client is nil")
	}

	if !strings.HasPrefix(baseURL, "http://") {
		baseURL = "http://" + baseURL
	}

	req := client.R().
		SetHeader("Content-Type", "text/plain").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(metricValue)

	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write([]byte(metricValue)); err != nil {
			return err
		}
		if err := gz.Close(); err != nil {
			return err
		}

		req.SetBody(b.Bytes())
		req.SetHeader("Content-Encoding", "gzip")
	}

	// fmt.Printf("Request Headers: TEXT %v\n", req.Header)

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/%s/%s/%s", baseURL, metricType, metricName, metricValue))
	if err != nil {
		logger.Log.WithError(err).Errorf("failed to send %s metric %s", metricType, metricName)
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
		logger.Log.Error(err)
		return err
	}

	return nil
}

func (j *JSONMetricSender) SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	var delta *int64
	var value *float64

	if metricType == "counter" {
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			logger.Log.Errorf("failed to parse counter value: %v", err)
			return fmt.Errorf("failed to parse counter value: %v", err)
		}
		delta = &val
	} else if metricType == "gauge" {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			logger.Log.Errorf("failed to parse gauge value: %v", err)
			return fmt.Errorf("failed to parse gauge value: %v", err)
		}
		value = &val
	} else {
		logger.Log.Errorf("invalid metric type: %s", metricType)
		return fmt.Errorf("invalid metric type: %s", metricType)
	}

	metric := []MetricsJSON{
		{
			ID:    metricName,
			MType: metricType,
			Delta: delta,
			Value: value,
		},
	}
	// add metrix slice to metrics
	// and on a handle update json, iterate over the slice

	if client == nil {
		return errors.New("client is nil")
	}

	if !strings.HasPrefix(baseURL, "http://") {
		baseURL = "http://" + baseURL
	}

	jsonData, err := json.Marshal(metric)
	if err != nil {
		err := fmt.Errorf("failed to marshal metric: %v", err)
		logger.Log.Error(err)
		return err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(jsonData); err != nil {
		err := fmt.Errorf("failed to write gzip data: %v", err)
		logger.Log.Error(err)
		return err
	}
	if err := gz.Close(); err != nil {
		err := fmt.Errorf("failed to close gzip writer: %v", err)
		logger.Log.Error(err)
		return err
	}

	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf.Bytes())

	resp, err := req.Post(fmt.Sprintf("%s/update/", baseURL))
	if err != nil {
		logger.Log.Errorf("failed to send %s metric %s: %v", metricType, metricName, err)
		return fmt.Errorf("failed to send %s metric %s: %v", metricType, metricName, err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Errorf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
		return fmt.Errorf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
	}

	return nil
}

func (j *JSONMetricSender) SendMetrics(client *resty.Client, metrics []Metric, baseURL string) error {
	for _, metric := range metrics {
		err := j.SendMetric(client, metric.Type, metric.Name, metric.Value, baseURL)
		if err != nil {
			return err
		}
	}
	return nil
}
