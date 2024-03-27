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
	GetServerAddress()
	GetStoreInterval()
	GetFileStoragePath()
	GetRestore()

	// Check that the flag values have been overridden
	if GetServerAddress() != "http://testaddress:8080" {
		t.Errorf("expected %v, got %v", "http://testaddress:8080", GetServerAddress())
	}
	if GetStoreInterval() != 200 {
		t.Errorf("expected %v, got %v", 200, GetStoreInterval())
	}
	if GetFileStoragePath() != "/tmp/test-metrics-db.json" {
		t.Errorf("expected %v, got %v", "/tmp/test-metrics-db.json", GetFileStoragePath())
	}
	if GetRestore() != false {
		t.Errorf("expected %v, got %v", false, GetRestore())
	}

	// Unset the environment variables to avoid affecting other tests
	os.Unsetenv("ADDRESS")
	os.Unsetenv("STORE_INTERVAL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("RESTORE")
}
