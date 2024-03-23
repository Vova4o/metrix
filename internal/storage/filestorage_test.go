package storage

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFileStorage(t *testing.T) {
	tests := []struct {
		name            string
		storeInterval   int
		fileStoragePath string
		restore         bool
		wantErr         bool
	}{
		{
			name:            "Valid parameters",
			storeInterval:   200,
			fileStoragePath: "/tmp/test-metrics-db.json",
			restore:         false,
			wantErr:         false,
		},
		{
			name:            "Invalid store interval",
			storeInterval:   0,
			fileStoragePath: "/tmp/test-metrics-db.json",
			restore:         false,
			wantErr:         true,
		},
		{
			name:            "Empty file storage path",
			storeInterval:   200,
			fileStoragePath: "",
			restore:         false,
			wantErr:         true,
		},
		{
			name:            "Restore with empty file storage path",
			storeInterval:   200,
			fileStoragePath: "",
			restore:         true,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := NewMemStorage()

			_, err := NewFileStorage(memStorage, tt.storeInterval, tt.fileStoragePath, tt.restore)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileStorage_LoadFromFile_SaveToFile(t *testing.T) {
	tests := []struct {
		name            string
		storeInterval   int
		fileStoragePath string
		restore         bool
		setupData       *MemStorage
		wantErr         bool
	}{
		{
			name:            "Valid data",
			storeInterval:   200,
			fileStoragePath: "/tmp/test-metrics-db.json",
			restore:         false,
			setupData: &MemStorage{
				GaugeMetrics: map[string]float64{
					"Alloc":       2139136,
					"BuckHashSys": 7708,
				},
				CounterMetrics: map[string]int64{
					"PollCount": 25,
				},
				Err: nil,
			},
			wantErr: false,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := NewMemStorage()

			fs, err := NewFileStorage(memStorage, tt.storeInterval, tt.fileStoragePath, tt.restore)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			jsonData, err := json.Marshal(tt.setupData)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			err = os.WriteFile(fs.fileStoragePath, jsonData, 0o644)
			if err != nil {
				t.Fatalf("failed to write data to file: %v", err)
			}

			err = fs.LoadFromFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Save to file
			err = fs.SaveToFile()
			if err != nil {
				t.Errorf("SaveToFile() error = %v", err)
			}

			// Load again to verify data was saved correctly
			err = fs.LoadFromFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() after SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileStorage(t *testing.T) {
	// Create a temporary file
	testMemStorage := MemStorage{
		GaugeMetrics: map[string]float64{
			"Alloc":       2139136,
			"BuckHashSys": 7708,
		},
		CounterMetrics: map[string]int64{
			"PollCount": 25,
		},
		Err: nil,
	}

	// created a test mem and file storage
	_, err := NewFileStorage(&testMemStorage, 1, "/tmp/test-metrics-db.json", true)
	if err != nil {
		// need assert function here
		assert.NoError(t, err)
	}

	// add timer to wait 2 seconds
	time.Sleep(2 * time.Second)

	// Create a new instance of FileStorage
	memStorage := NewMemStorage()
	fs1, err := NewFileStorage(memStorage, 1, "/tmp/test-metrics-db.json", true)
	if err != nil {
		t.Fatal(err)
	}

	// The values in fs1 should be the same as the ones in testMemStorage
	gauge1, _ := fs1.GetGauge("gauge")
	counter1, _ := fs1.GetCounter("counter")
	gaugeTest, _ := testMemStorage.GetGauge("gauge")
	counterTest, _ := testMemStorage.GetCounter("counter")
	assert.Equal(t, gaugeTest, gauge1)
	assert.Equal(t, counterTest, counter1)
}
