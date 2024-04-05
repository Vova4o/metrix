package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MetricValue(s Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")

		value, err := getMetricValue(s, metricType, metricName)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, fmt.Sprintf("%v", value))
	}
}

func MetricValueJSON(s Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metrics MetricsJSON

		// Bind the JSON request body into the metrics struct
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		value, err := getJSONValue(s, metrics)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if metrics.MType == "gauge" {
			c.JSON(http.StatusOK, gin.H{
				"id":    metrics.ID,
				"type":  metrics.MType,
				"value": value,
			})
		} else if metrics.MType == "counter" {
			c.JSON(http.StatusOK, gin.H{
				"id":    metrics.ID,
				"type":  metrics.MType,
				"delta": value,
			})
		}
	}
}
