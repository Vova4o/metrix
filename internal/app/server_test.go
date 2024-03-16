package app

import (
    "net"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
    // Start the server in a goroutine to allow the test to continue and make requests
    go NewServer()

    // Allow some time for the server to start
    time.Sleep(time.Second)

    // Check if the server is up and running
    _, err := net.Dial("tcp", "localhost:8080")
    assert.NoError(t, err, "server did not start")
}