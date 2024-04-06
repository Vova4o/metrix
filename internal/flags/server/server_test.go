package serverflags

import (
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Set environment variables to override the default flag values
	os.Setenv("ADDRESS", "http://testaddress:8080")
	os.Setenv("STORE_INTERVAL", "200")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/test-metrics-db.json")
	os.Setenv("RESTORE", "false")

	// Call the function
	ServerAddress()
	StoreInterval()
	FileStoragePath()
	Restore()

	// Check that the flag values have been overridden
	if ServerAddress() != "http://testaddress:8080" {
		t.Errorf("expected %v, got %v", "http://testaddress:8080", ServerAddress())
	}
	if StoreInterval() != 200 {
		t.Errorf("expected %v, got %v", 200, StoreInterval())
	}
	if FileStoragePath() != "/tmp/test-metrics-db.json" {
		t.Errorf("expected %v, got %v", "/tmp/test-metrics-db.json", FileStoragePath())
	}
	if Restore() != false {
		t.Errorf("expected %v, got %v", false, Restore())
	}

	// Unset the environment variables to avoid affecting other tests
	os.Unsetenv("ADDRESS")
	os.Unsetenv("STORE_INTERVAL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("RESTORE")
}
