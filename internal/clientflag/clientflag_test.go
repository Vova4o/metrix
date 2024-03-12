package clientflag

import (
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Set environment variables to override the default flag values
	os.Setenv("ADDRESS", "http://testaddress:8080")
	os.Setenv("REPORT_INTERVAL", "20")
	os.Setenv("POLL_INTERVAL", "5")
	os.Setenv("SERVER_ADDRESS", "testaddress:8080")

	// Call the function
	ParseFlags()

	// Check that the flag values have been overridden
	if *ServerAddress != "http://testaddress:8080" {
		t.Errorf("expected %v, got %v", "http://testaddress:8080", *ServerAddress)
	}
	if *ReportInterval != 20 {
		t.Errorf("expected %v, got %v", 20, *ReportInterval)
	}
	if *PollInterval != 5 {
		t.Errorf("expected %v, got %v", 5, *PollInterval)
	}

	// Unset the environment variables to avoid affecting other tests
	os.Unsetenv("ADDRESS")
	os.Unsetenv("REPORT_INTERVAL")
	os.Unsetenv("POLL_INTERVAL")
	os.Unsetenv("SERVER_ADDRESS")
}
