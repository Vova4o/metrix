package config

import (
	"os"
)

const (
	AgentLogFile  = "agentlog.txt"
)

var (
	LogfileAgent  *os.File
)
