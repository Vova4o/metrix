package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MetricValue(s Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")

		var mt Metricer
		switch metricType {
		case "gauge":
			mt = GaugeMetricType{}
		case "counter":
			mt = CounterMetricType{}
		default:
			log.Printf("Invalid metric type: %s", metricType)
			c.String(http.StatusBadRequest, "Invalid metric type")
			return
		}

		value, exists := mt.GetValue(s, metricName)
		if !exists {
			c.String(http.StatusNotFound, "Metric not found")
			return
		}

		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, mt.FormatValue(value))
	}
}

func MetricValueJSON(s Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metrics MetricsJSON

		// Bind the JSON request body into the metrics struct
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var mt Metricer
		switch metrics.MType {
		case "gauge":
			mt = GaugeMetricType{}
		case "counter":
			mt = CounterMetricType{}
		default:
			log.Printf("Invalid metric type: %s", metrics.MType)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
			return
		}

		value, exists := mt.GetValue(s, metrics.ID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			return
		}

		if metrics.MType == "gauge" {
			c.JSON(http.StatusOK, gin.H{
				"id":    metrics.ID,
				"type":  metrics.MType,
				"value": value,
			})
			return
		}
		if metrics.MType == "counter" {
			c.JSON(http.StatusOK, gin.H{
				"id":    metrics.ID,
				"type":  metrics.MType,
				"delta": value,
			})
		}
	}
}
