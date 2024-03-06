package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ServerAddress  = flag.String("address", "localhost:8080", "HTTP server network address")
	ReportInterval = flag.Duration("report_interval", 10*time.Second, "Interval between fetching reportable metrics")
	PollInterval   = flag.Duration("poll_interval", 2*time.Second, "Interval between polling metrics")
)

func parseFlags() {
	flag.Parse()

	// Override the server address, report interval, and poll interval with environment variables if they are set
	if address := os.Getenv("ADDRESS"); address != "" {
		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			address = "http://" + address
		}
		*ServerAddress = address
	}

	if ri := os.Getenv("REPORT_INTERVAL"); ri != "" {
		if riInt, err := strconv.Atoi(ri); err == nil {
			*ReportInterval = time.Duration(riInt) * time.Second
		}
	}

	if pi := os.Getenv("POLL_INTERVAL"); pi != "" {
		if piInt, err := strconv.Atoi(pi); err == nil {
			*PollInterval = time.Duration(piInt) * time.Second
		}
	}
}
