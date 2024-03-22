package clientmetrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"Vova4o/metrix/internal/logger"

	"github.com/go-resty/resty/v2"
)

func (t *TextMetricSender) SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error {
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

	fmt.Printf("Request Headers: TEXT %v\n", req.Header)

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/%s/%s/%s", baseURL, metricType, metricName, metricValue))
	if err != nil {
		logger.Log.Logger.WithError(err).Errorf("failed to send %s metric %s", metricType, metricName)
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
		logger.Log.Logger.Error(err)
		return err
	}

	return nil
}

func (j *JSONMetricSender) SendMetric(client *resty.Client, metricType, metricName, metricValue, baseURL string) error {
	var delta *int64
	var value *float64

	if metricType == "counter" {
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			logger.Log.Logger.Errorf("failed to parse counter value: %v", err)
			return fmt.Errorf("failed to parse counter value: %v", err)
		}
		delta = &val
	} else if metricType == "gauge" {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			logger.Log.Logger.Errorf("failed to parse gauge value: %v", err)
			return fmt.Errorf("failed to parse gauge value: %v", err)
		}
		value = &val
	} else {
		logger.Log.Logger.Errorf("invalid metric type: %s", metricType)
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
	// this is how the error idealy should look like.
	if err != nil {
		err := fmt.Errorf("failed to marshal metric: %v", err)
		logger.Log.Logger.Error(err)
		return err
	}

	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(jsonDate)

	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write(jsonDate); err != nil {
			return err
		}
		if err := gz.Close(); err != nil {
			return err
		}

		req.SetBody(b.Bytes())
		req.SetHeader("Content-Encoding", "gzip")
	} else {
		req.SetBody(jsonDate)
	}

	fmt.Printf("Request Headers: JSON %v\n", req.Header)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Post(fmt.Sprintf("%s/update/", baseURL))
		// this is how it maigh look like
	if err != nil {
		logger.Log.Logger.Errorf("failed to send %s metric %s: %v", metricType, metricName, err)
		return fmt.Errorf("failed to send %s metric %s: %v", metricType, metricName, err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Logger.Errorf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
		return fmt.Errorf("server returned non-OK status for %s metric %s: %v", metricType, metricName, resp.Status())
	}

	return nil
}

// GzipWriter is a middleware that compresses response body in gzip format
func GzipWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (gzw gzipResponseWriter) Write(b []byte) (int, error) {
	return gzw.Writer.Write(b)
}
