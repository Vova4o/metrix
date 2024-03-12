package clientmetrics

import (
	"Vova4o/metrix/internal/config"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
)

// HttpClient is an interface for making HTTP requests
type RestClient interface {
	R() *resty.Request
}

// reportMetrics sends the metrics to the server
func ReportMetrics(baseURL string) {
	client := resty.New()

	// Add a middleware logger
	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		log.Printf("Sending request to URL: %s\n", request.URL)
		return nil
	})

	client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		log.Printf("Received response with status code: %d\n", response.StatusCode())
		return nil
	})

	// Use the client to send a request
	resp, err := client.R().Get("http://example.com")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	log.Printf("Response: %s", resp)

	errs := make(chan error)
	var wg sync.WaitGroup

	for metricName, metricValue := range config.CounterMetrics {
		wg.Add(1)
		go func(name string, value float64) {
			defer wg.Done()
			err := SendMetric(client, "gauge", name, value, baseURL)
			if err != nil {
				log.Printf("error sending gauge metric %s: %v", name, err)
				errs <- fmt.Errorf("error sending gauge metric %s: %v", name, err)
			}
		}(metricName, float64(metricValue))
	}

	for metricName, metricValue := range config.CounterMetrics {
		wg.Add(1)
		go func(name string, value float64) {
			defer wg.Done()
			err := SendMetric(client, "counter", name, value, baseURL)
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
