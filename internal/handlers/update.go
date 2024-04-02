package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUpdateText(s Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")
		metricValue := c.Param("metricValue")

		var mt Metricer
		switch metricType {
		case "gauge":
			mt = GaugeMetricType{}
		case "counter":
			mt = CounterMetricType{}
		default:
			log.Printf("Invalid metric type: %s", metricType)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
			return
		}

		value, err := mt.ParseValue(metricValue)
		if err != nil {
			log.Printf("Invalid metric value: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value"})
			return
		}

		mt.Store(s, metricName, value)

		c.Status(http.StatusOK)
	}
}

func HandleUpdateJSON(s Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metrics MetricsJSON
		err := c.ShouldBindJSON(&metrics)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if metrics.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id"})
			return
		}

		var mt Metricer
		var value interface{}
		switch metrics.MType {
		case "gauge":
			mt = GaugeMetricType{}
			if metrics.Value != nil {
				value = *metrics.Value
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Value is required for gauge type"})
				return
			}
		case "counter":
			mt = CounterMetricType{}
			if metrics.Delta != nil {
				value = *metrics.Delta
			}
		default:
			log.Printf("Invalid metric type: %s", metrics.MType)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
			return
		}

		mt.Store(s, metrics.ID, value)

		// Get the latest value from the storage
		latestValue, ok := mt.GetValue(s, metrics.ID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest value"})
			return
		}

		// Update the metrics value based on the type
		if metrics.MType == "counter" {
			if val, ok := latestValue.(int64); ok {
				metrics.Delta = &val
			} else {
				log.Printf("Expected *int64, got %T", latestValue)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		} else if metrics.MType == "gauge" {
			if val, ok := latestValue.(float64); ok {
				metrics.Value = &val
			} else {
				log.Printf("Expected *float64, got %T", latestValue)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		}

		c.JSON(http.StatusOK, metrics)
	}
}
