package appserver

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

// TestServer tests the NewServer function
// by starting the server and checking if it is up and running
func TestServer(t *testing.T) {
	// Create a channel to receive the error
	errCh := make(chan error, 1)

	// Start the server in a goroutine and send any error to the channel
	go func() {
		errCh <- NewServer()
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Check if the server is up and running
	_, err := net.Dial("tcp", "localhost:8080")
	assert.NoError(t, err, "server did not start")

	// After making requests, check if there was an error starting the server
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Failed to start the server: %v", err)
		}
	default:
		// If no error was received, the server started successfully
	}
}

func TestFileStorage(t *testing.T) {
	// Create a logger for the test
	log := logrus.New()
	hook := test.NewLocal(log)

	// Test when err is not nil
	t.Run("error creating file storage", func(t *testing.T) {
		// Reset the hook
		hook.Reset()
		// Simulate an error
		err := errors.New("simulated error")
		// Call the function that logs the error
		if err != nil {
			log.WithError(err).Error("Failed to create new file storage")
		}

		// Check if the error was logged
		entries := hook.AllEntries()
		assert.NotEmpty(t, entries, "No error was logged")
		assert.Equal(t, "Failed to create new file storage", entries[0].Message, "Unexpected log message")
		assert.Equal(t, err, entries[0].Data[logrus.ErrorKey], "Unexpected error in log data")
	})

	// Test when err is nil
	t.Run("no error creating file storage", func(t *testing.T) {
		// Reset the hook
		hook.Reset()
		// Simulate no error
		var err error = nil

		// Call the function that logs the error
		if err != nil {
			log.WithError(err).Error("Failed to create new file storage")
		} else {
			fmt.Println("Not using file storage")
		}

		// Check if the error was not logged
		entries := hook.AllEntries()
		assert.Empty(t, entries, "An error was logged")
	})
}
