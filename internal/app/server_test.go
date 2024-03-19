package app

import (
    "net"
    "testing"
    "time"

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