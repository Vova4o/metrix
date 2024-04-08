package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestShowMetrics(t *testing.T) {
	// Create a mock storager
	s := &mockStorager{}

	// Create a test Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set the request and response
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	c.Request = req

	// Call the ShowMetrics handler
	ShowMetrics(s, "metrix.page.tmpl")(c)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, w.Code)
	}

	// Check if the response body contains certain fields
	body := w.Body.String()
	if !strings.Contains(body, "<h1>Gauge Metrics</h1>") {
		t.Errorf("expected body to contain %v, body was %v", "<h1>Gauge Metrics</h1>", body)
	}
	if !strings.Contains(body, "<h1>Counter Metrics</h1>") {
		t.Errorf("expected body to contain %v, body was %v", "<h1>Counter Metrics</h1>", body)
	}
}
