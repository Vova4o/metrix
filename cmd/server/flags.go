package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var serverAddress = flag.String("a", "localhost:8080", "HTTP server address")

func parseFlags() {
	// Parse the flags
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		// If there was an error, print it and exit with a non-zero status code
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Check if serverAddress starts with http://
	if strings.HasPrefix(*serverAddress, "http://") {
		// Remove http:// from serverAddress
		*serverAddress = strings.TrimPrefix(*serverAddress, "http://")
	}
}
