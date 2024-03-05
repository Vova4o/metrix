package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHttpClient is a mock implementation of the HttpClient interface
type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestSendMetric(t *testing.T) {
    mockClient := new(MockHttpClient)
    client := HttpClient(mockClient)

    t.Run("Successful metric send", func(t *testing.T) {
        mockClient.On("Post", "http://localhost:8080/update/gauge/testGauge/42", "text/plain", strings.NewReader("")).Return(&http.Response{
            StatusCode: 200,
            Body:       io.NopCloser(bytes.NewBufferString("OK")),
        }, nil)
        err := sendMetric(client, "gauge", "testGauge", 42)
        assert.NoError(t, err)
        mockClient.ExpectedCalls = []*mock.Call{}
    })

    t.Run("Failed metric send", func(t *testing.T) {
        mockClient.On("Post", "http://localhost:8080/update/gauge/testGauge/42", "text/plain", strings.NewReader("")).Return(&http.Response{
            StatusCode: 500,
            Body:       io.NopCloser(bytes.NewBufferString("")),
        }, errors.New("failed to send metric"))
        err := sendMetric(client, "gauge", "testGauge", 42)
        assert.Error(t, err)
        mockClient.ExpectedCalls = []*mock.Call{}
    })

    t.Run("Non-OK response status", func(t *testing.T) {
        mockClient.On("Post", "http://localhost:8080/update/gauge/testGauge/42", "text/plain", strings.NewReader("")).Return(&http.Response{
            StatusCode: 500,
            Body:       io.NopCloser(bytes.NewBufferString("")),
        }, nil)
        err := sendMetric(client, "gauge", "testGauge", 42)
        assert.Error(t, err)
        mockClient.ExpectedCalls = []*mock.Call{}
    })
}
