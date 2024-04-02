package middleware

import (
	"compress/gzip"
	"strings"
	"time"

	"Vova4o/metrix/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type gzipWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		logger.Log.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"duration": duration.String(),
			"size":     c.Writer.Size(),
		}).Info("Handled request")
	}
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipMiddleware compresses response body in gzip format if the client supports it
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the client accepts gzip compression
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		c.Writer.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		c.Writer = &gzipWriter{ResponseWriter: c.Writer, Writer: gz}
		c.Next()
	}
}
