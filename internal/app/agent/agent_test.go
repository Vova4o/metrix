package appagent_test

import (
	"context"
	"testing"
	"time"

	appagent "Vova4o/metrix/internal/app/agent"
	"Vova4o/metrix/internal/logger"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewAgent(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		messages []string
	}{
		{
			name:     "Test URL",
			url:      "http://example.com",
			messages: []string{"Sending request", "Received response"},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new resty client
			client := resty.New()

			err := logger.New("test.log")
			if err != nil {
				logger.Log.WithError(err).Error("Failed to initialize logger")
			}

			defer func() {
				if err := logger.Close(); err != nil {
					t.Errorf("Failed to close log file: %v\n", err)
				}
			}()

			// Create a test hook for the logger
			hook := test.NewLocal(logger.Log)

			// Create a context with a timeout
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// Run the NewAgent function with the context as the first argument
			err = appagent.NewAgent(ctx, client)
			if err != nil {
				t.Fatalf("Failed to start the agent: %v", err)
			}

			// Make a request
			_, _ = client.R().Get(tt.url)

			// Check if the request was logged
			assert.Equal(t, len(tt.messages), len(hook.Entries))
			for i, message := range tt.messages {
				assert.Equal(t, message, hook.Entries[i].Message)
				if i == 0 { // Only the first log entry should have the URL
					assert.Equal(t, tt.url, hook.Entries[i].Data["url"])
				}
			}
		})
	}
}
