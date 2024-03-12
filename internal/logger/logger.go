package logger

import (
	"fmt"
	"os"
)

func Logger(name string) (*os.File, error) {
	logFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return nil, err
	}
	// defer logFile.Close()

	return logFile, nil
}
