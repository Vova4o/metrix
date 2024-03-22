package config

import (
	"os"
)

const (
	ServerLogFile = "serverlog.txt"
)

var (
	LogfileServer *os.File
)
