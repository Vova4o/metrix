package middleware

import (
	"compress/gzip"
	"net/http"
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

// RequestLogger logs all the requests
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
func CompressGzip() gin.HandlerFunc {
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

// OptionalDecompressGzip decompresses request body if it is compressed in gzip format
// so you dont need to handle this gzip in your handles
func DecompressGzip(c *gin.Context) {
	if c.GetHeader("Content-Encoding") != "gzip" {
		c.Next()
		return
	}

	r, err := gzip.NewReader(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer r.Close()

	c.Request.Header.Del("Content-Encoding")
	c.Request.Header.Del("Content-Length")
	c.Request.Body = r

	c.Next()
}
