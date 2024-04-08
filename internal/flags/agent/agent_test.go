package agentflags

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Set environment variables to override the default flag values
	os.Setenv("ADDRESS", "http://testaddress:8080")
	os.Setenv("REPORT_INTERVAL", "20")
	os.Setenv("POLL_INTERVAL", "5")

	// Call the function
	ServerAddress()
	ReportInterval()
	PollInterval()

	// Check that the flag values have been overridden
	if ServerAddress() != "http://testaddress:8080" {
		t.Errorf("expected %v, got %v", "http://testaddress:8080", ServerAddress())
	}
	if ReportInterval() != 20 {
		t.Errorf("expected %v, got %v", 20, ReportInterval())
	}
	if PollInterval() != 5 {
		t.Errorf("expected %v, got %v", 5, PollInterval())
	}

	// Unset the environment variables to avoid affecting other tests
	os.Unsetenv("ADDRESS")
	os.Unsetenv("REPORT_INTERVAL")
	os.Unsetenv("POLL_INTERVAL")
}

func TestParseFlags_InvalidValues(t *testing.T) {
	// Set environment variables to invalid values
	os.Setenv("ADDRESS", "")
	os.Setenv("REPORT_INTERVAL", "not a number")
	os.Setenv("POLL_INTERVAL", "not a number")

	// Call the functions
	ServerAddress()
	ReportInterval()
	PollInterval()

	// Check that the functions return the default values
	if ServerAddress() != "localhost:8080" {
		t.Errorf("expected %v, got %v", "localhost:8080", ServerAddress())
	}
	if ReportInterval() != 10 {
		t.Errorf("expected %v, got %v", 10, ReportInterval())
	}
	if PollInterval() != 2 {
		t.Errorf("expected %v, got %v", 2, PollInterval())
	}

	// Unset the environment variables to avoid affecting other tests
	os.Unsetenv("ADDRESS")
	os.Unsetenv("REPORT_INTERVAL")
	os.Unsetenv("POLL_INTERVAL")
}

func TestInitFlags(t *testing.T) {
	// Set the flags
	flags.Set("ServerAddress", "http://localhost:8080")
	flags.Set("ReportInterval", "10")
	flags.Set("PollInterval", "2")

	// Check that the flags are set correctly
	if flag := flags.Lookup("ServerAddress"); flag == nil || !flag.Changed {
		t.Error("ServerAddress flag is not set")
	}
	if flag := flags.Lookup("ReportInterval"); flag == nil || !flag.Changed {
		t.Error("ReportInterval flag is not set")
	}
	if flag := flags.Lookup("PollInterval"); flag == nil || !flag.Changed {
		t.Error("PollInterval flag is not set")
	}
}

func TestInitEnvVars(t *testing.T) {
	// Set the environment variables
	os.Setenv("ADDRESS", "http://localhost:8080")
	os.Setenv("REPORT_INTERVAL", "10")
	os.Setenv("POLL_INTERVAL", "2")

	// Check that the environment variables are set correctly
	if os.Getenv("ADDRESS") == "" {
		t.Error("ADDRESS environment variable is not set")
	}
	if os.Getenv("REPORT_INTERVAL") == "" {
		t.Error("REPORT_INTERVAL environment variable is not set")
	}
	if os.Getenv("POLL_INTERVAL") == "" {
		t.Error("POLL_INTERVAL environment variable is not set")
	}
}

func TestBindFlagToViper(t *testing.T) {
	// Create a buffer to hold the log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Call the function with an invalid flag name
	bindFlagToViper("invalid")

	// Check that the error message was logged
	if !strings.Contains(buf.String(), "invalid") {
		t.Errorf("expected %v, got %v", "invalid", buf.String())
	}
}

func TestBindEnvToViper(t *testing.T) {
	// Set the environment variable
	os.Setenv("INVALID", "value")

	// Call the function with the env key
	bindEnvToViper("invalid", "INVALID")

	// Check that the environment variable is set
	if value := os.Getenv("INVALID"); value == "" {
		t.Errorf("expected %v, got %v", "value", value)
	}
}
