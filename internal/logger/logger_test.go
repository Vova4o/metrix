package logger

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Create a temporary log file for testing
	logFile := "test.log"
	defer os.Remove(logFile)

	err := New(logFile)
	assert.NoError(t, err)

	// Check if the log file is created
	_, err = os.Stat(logFile)
	assert.NoError(t, err)

	// Open the log file
	file, err := os.Open(logFile)
	assert.NoError(t, err)
	defer file.Close()

	// Check if the logrus logger is properly configured
	assert.IsType(t, &logrus.Logger{}, Log)
	assert.Equal(t, file.Name(), logFile)
	assert.IsType(t, &logrus.JSONFormatter{}, Log.Formatter)
}
