package main

import (
	"flag"
	"os"
	"strings"
)

var serverAddress = flag.String("a", "", "HTTP server address")

func parseFlags() {
	// Parse the flags
	flag.Parse()

	// Check if serverAddress was provided as a flag
	if *serverAddress == "" {
		// If not, check if the ADDRESS environment variable is set
		envAddress := os.Getenv("ADDRESS")
		if envAddress != "" {
			*serverAddress = envAddress
		} else {
			// If not, use the default value
			*serverAddress = "localhost:8080"
		}
	}

	// Check if serverAddress starts with http://
	if strings.HasPrefix(*serverAddress, "http://") {
		// Remove http:// from serverAddress
		*serverAddress = strings.TrimPrefix(*serverAddress, "http://")
	}
}
