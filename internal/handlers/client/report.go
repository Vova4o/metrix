package clientmetrics

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"Vova4o/metrix/internal/config"
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
		logrus.WithFields(logrus.Fields{
			"url": request.URL,
		}).Info("Sending request")

		return nil
	})

	client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		logrus.WithFields(logrus.Fields{
			"status": response.StatusCode(),
			"body":   response.String(),
		}).Info("Received response")

		return nil
	})

	errs := make(chan error)
	var wg sync.WaitGroup

	for metricName, metricValue := range config.GaugeMetrics {
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
		go func(name string, value int64) {
			defer wg.Done()
			err := SendMetric(client, "counter", name, float64(value), baseURL)
			if err != nil {
				log.Printf("error sending counter metric %s: %v", name, err)
				errs <- fmt.Errorf("error sending counter metric %s: %v", name, err)
			}
		}(metricName, metricValue)
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
