package config

import (
	"os"
)

const (
	ServerLogFile = "serverlog.txt"
	AgentLogFile  = "agentlog.txt"
)

var (
	LogfileServer *os.File
	LogfileAgent  *os.File
)

// // Variables to store the command-line flags
// var (
// 	ServerAddress  = flag.String("a", "localhost:8080", "HTTP server network address")
// 	ReportInterval = flag.Int("r", 10, "Interval between fetching reportable metrics in seconds")
// 	PollInterval   = flag.Int("p", 2, "Interval between polling metrics in seconds")
// )
