package main

import (
	"flag"
	"os"
	"strings"
)

// parseFlags parses the flags and sets the serverAddress variable
var ServerAddress = flag.String("a", "", "HTTP server address")

func parseFlags() string {
	// Parse the flags
	flag.Parse()

	// Check if serverAddress was provided as a flag
	if *ServerAddress == "" {
		// If not, check if the ADDRESS environment variable is set
		envAddress := os.Getenv("ADDRESS")
		// If it is, use the value from the environment variable
		if envAddress != "" {
			*ServerAddress = envAddress
		} else {
			// If not, use the default value
			*ServerAddress = "localhost:8080"
		}
	}

	// Check if serverAddress starts with http://
	if strings.HasPrefix(*ServerAddress, "http://") {
		// Remove http:// from serverAddress
		*ServerAddress = strings.TrimPrefix(*ServerAddress, "http://")
	}
	return *ServerAddress
}

// func parseFlags() {
// 	// Parse the flags
// 	flag.Parse()

// 	// Check if serverAddress was provided as a flag
// 	if *ServerAddress == "" {
// 		// If not, check if the ADDRESS environment variable is set
// 		envAddress := os.Getenv("ADDRESS")
// 		// If it is, use the value from the environment variable
// 		if envAddress != "" {
// 			*ServerAddress = envAddress
// 		} else {
// 			// If not, use the default value
// 			*ServerAddress = "localhost:8080"
// 		}
// 	}

// 	// Check if serverAddress starts with http://
// 	if strings.HasPrefix(*ServerAddress, "http://") {
// 		// Remove http:// from serverAddress
// 		*ServerAddress = strings.TrimPrefix(*ServerAddress, "http://")
// 	}
// }
