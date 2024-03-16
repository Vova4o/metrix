package clientmetrics

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func SendMetric(client RestClient, metricType, metricName string, metricValue float64, baseURL string) error {
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
