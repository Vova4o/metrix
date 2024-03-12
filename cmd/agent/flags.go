package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

// Variables to store the command-line flags
var (
	ServerAddress  = flag.String("a", "localhost:8080", "HTTP server network address")
	ReportInterval = flag.Int("r", 10, "Interval between fetching reportable metrics in seconds")
	PollInterval   = flag.Int("p", 2, "Interval between polling metrics in seconds")
)

func parseFlags() {
	// Parse the command-line flags
	flag.Parse()

	// Override the server address, report interval, and poll interval with environment variables if they are set
	if address := os.Getenv("ADDRESS"); address != "" {
		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			address = "http://" + address
		}
		*ServerAddress = address
	}

	// If the REPORT_INTERVAL environment variable is set, override the default value
	if ri := os.Getenv("REPORT_INTERVAL"); ri != "" {
		if riInt, err := strconv.Atoi(ri); err == nil {
			*ReportInterval = riInt
		}
	}

	// If the POLL_INTERVAL environment variable is set, override the default value
	if pi := os.Getenv("POLL_INTERVAL"); pi != "" {
		if piInt, err := strconv.Atoi(pi); err == nil {
			*PollInterval = piInt
		}
	}
}
