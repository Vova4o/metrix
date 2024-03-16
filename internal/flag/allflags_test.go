package allflags

import (
    "os"
    "testing"
)

func TestParseFlags(t *testing.T) {
    // Set environment variables to override the default flag values
    os.Setenv("ADDRESS", "http://testaddress:8080")
    os.Setenv("REPORT_INTERVAL", "20")
    os.Setenv("POLL_INTERVAL", "5")

    // Call the function
    GetServerAddress()
    GetReportInterval()
    GetPollInterval()

    // Check that the flag values have been overridden
    if GetServerAddress() != "http://testaddress:8080" {
        t.Errorf("expected %v, got %v", "http://testaddress:8080", GetServerAddress())
    }
    if GetReportInterval() != 20 {
        t.Errorf("expected %v, got %v", 20, GetReportInterval())
    }
    if GetPollInterval() != 5 {
        t.Errorf("expected %v, got %v", 5, GetPollInterval())
    }

    // Unset the environment variables to avoid affecting other tests
    os.Unsetenv("ADDRESS")
    os.Unsetenv("REPORT_INTERVAL")
    os.Unsetenv("POLL_INTERVAL")
}