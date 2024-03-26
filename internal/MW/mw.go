package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"time"

	"Vova4o/metrix/internal/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func RequestLogger(logger *logger.FileLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()

			defer func() {
				duration := time.Since(start)
				logger.Logger.WithFields(logrus.Fields{
					"status":   ww.Status(),
					"method":   r.Method,
					"path":     r.URL.Path,
					"duration": duration.String(),
					"size":     ww.BytesWritten(),
				}).Info("Handled request")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// GzipMiddleware compresses response body in gzip format if the client supports it
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client accepts gzip compression
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
