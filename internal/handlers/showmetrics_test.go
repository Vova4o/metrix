package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/storage"
)
func TestShowMetricsHandler(t *testing.T) {
    tests := []struct {
        name           string
        templateFile   string
        expectedStatus int
        expectedBody   []string
    }{
        {
            name:           "Success",
            templateFile:   "metrix.page.tmpl",
            expectedStatus: http.StatusOK,
            expectedBody:   []string{"<li>gaugeTest: 10.0000</li>", "<li>counterTest: 20</li>"},
        },
        {
            name:           "TemplateParseError",
            templateFile:   "nonexistentfile",
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   []string{"Internal Server Error\n"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create a storage and set some metrics
            storage := storage.NewMemStorage()
            storage.SetGauge("gaugeTest", 10.0000)
            storage.SetCounter("counterTest", 20)

            // Create a request to pass to our handler
            req, err := http.NewRequest("GET", "", nil)
            if err != nil {
                t.Fatal(err)
            }

            // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
            rr := httptest.NewRecorder()
            handler := handlers.ShowMetrics(storage, tt.templateFile)

            // Our handlers satisfy http.Handler, so we can call their ServeHTTP method
            // directly and pass in our Request and ResponseRecorder
            handler.ServeHTTP(rr, req)

            // Check the status code is what we expect
            assert.Equal(t, tt.expectedStatus, rr.Code)

            // Check the response body contains the expected metrics and their values
            body := rr.Body.String()
            for _, b := range tt.expectedBody {
                assert.Contains(t, body, b)
            }
        })
    }
}