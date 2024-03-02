package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestUpdateHandler(t *testing.T) {
    // Create a request to pass to our handler
    req := httptest.NewRequest("POST", "/update/gauge/test/10", nil)

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Create a HTTP handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // The same code as your handler
    })

    // Call ServeHTTP method directly and pass in our Request and ResponseRecorder
    handler.ServeHTTP(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

    // Check the response body
    expected := "" // Expected response body
    assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body")
}