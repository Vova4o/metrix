package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	ServerAddress  = flag.String("a", "http://localhost:8080", "HTTP server address")
	ReportInterval = flag.Duration("r", 10*time.Second, "report interval")
	PollInterval   = flag.Duration("p", 2*time.Second, "poll interval")
)

func parseFlags() {
	// Parse the flags
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		// If there was an error, print it and exit with a non-zero status code
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
