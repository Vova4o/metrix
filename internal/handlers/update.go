package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUpdateText(s Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")
		metricValue := c.Param("metricValue")

		// checking for empty values
		if metricType == "" || metricName == "" || metricValue == "" {
			log.Printf("Empty values")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty values"})
			return
		}

		err := storeMetric(s, metricType, metricName, metricValue)
		if err != nil {
			log.Printf("Failed to store metric: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

func HandleUpdateJSON(s Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metrics []MetricsJSON
		var singleMetric MetricsJSON
		var body []byte
		var err error

		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		fmt.Println(string(body))

		err = json.Unmarshal(body, &metrics)
		if err != nil {
			// If unmarshalling into an array fails, try unmarshalling into a single object
			err = json.Unmarshal(body, &singleMetric)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
				return
			}

			// If unmarshalling into a single object succeeds, add it to the metrics array
			metrics = append(metrics, singleMetric)
		}

		for _, metric := range metrics {
			if metric.ID == "" {
				log.Printf("Missing id")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id"})
				return
			}

			err := storeMetricJSON(s, metric)
			if err != nil {
				log.Printf("Failed to store metric: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		// Respond with a single JSON object if the metrics array contains only one element
		if len(metrics) > 1 {
			c.JSON(http.StatusOK, metrics)
		} else {
			c.JSON(http.StatusOK, metrics[0])
		}
	}
}
