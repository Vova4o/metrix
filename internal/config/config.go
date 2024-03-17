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
