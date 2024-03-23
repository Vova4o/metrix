package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"Vova4o/metrix/internal/logger"
)

type FileStorage struct {
	memStorage      *MemStorage
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func NewFileStorage(memStorage *MemStorage, storeInterval int, fileStoragePath string, restore bool) (*FileStorage, error) {
	fs := &FileStorage{
		memStorage:      memStorage,
		storeInterval:   storeInterval,
		fileStoragePath: fileStoragePath,
		restore:         restore,
	}

	if fs.storeInterval <= 0 {
		return nil, fmt.Errorf("storeInterval must be greater than 0")
	}
	if fs.fileStoragePath == "" {
		return nil, fmt.Errorf("fileStoragePath cannot be empty")
	}
	if fs.restore && fs.fileStoragePath == "" {
		return nil, fmt.Errorf("restore cannot be true if fileStoragePath is empty")
	}
	if fs.restore && fs.fileStoragePath != "" {
		if _, err := os.Stat(fs.fileStoragePath); os.IsNotExist(err) {
			// create a new file if it doesn't exist
			file, err := os.Create(fs.fileStoragePath)
			if err != nil {
				logger.Log.Logger.WithError(err).Error("Failed to create new file")
				return nil, err
			}
			defer file.Close()
			logger.Log.Logger.Info("No previous metrics file found. Created a new one.")
		} else if err != nil {
			// some other error occurred when trying to stat the file
			logger.Log.Logger.WithError(err).Error("Failed to check if metrics file exists")
			return nil, err
		}
	}

	// var fileExists bool
	// // Load previously saved metrics from the file at startup
	// if storage.restore && storage.fileStoragePath != "" {
	// 	if _, err := os.Stat(storage.fileStoragePath); os.IsNotExist(err) {
	// 		fileExists = false
	// 		// File does not exist, create a new one
	// 		file, err := os.Create(storage.fileStoragePath)
	// 		if err != nil {
	// 			logger.Log.Logger.WithError(err).Error("Failed to create new file")
	// 			return nil
	// 		}
	// 		defer file.Close()
	// 		logger.Log.Logger.Info("No previous metrics file found. Created a new one.")
	// 	} else if err != nil {
	// 		// Some other error occurred when trying to stat the file
	// 		logger.Log.Logger.WithError(err).Error("Failed to check if metrics file exists")
	// 		return nil
	// 	} else {
	// 		fileExists = true // Set fileExists to true when the file does exist
	// 	}

	// 	if fileExists {
	// 		// Check if storage is not nil before calling LoadFromFile
	// 		if err := storage.LoadFromFile(); err != nil {
	// 			logger.Log.Logger.WithError(err).Error("Failed to load metrics from file")
	// 		}
	// 	}
	// }

	// Save current metrics to the file at the specified interval
	if fs.fileStoragePath != "" {
		go fs.saveAtInterval()
	}

	return fs, nil
}

func (s *FileStorage) LoadFromFile() error {
	file, err := os.ReadFile(s.fileStoragePath)
	if err != nil {
		return err
	}

	// fmt.Printf("File contents: %s\n", string(file))

	err = json.Unmarshal(file, s.memStorage)
	if err != nil {
		fmt.Printf("Failed to unmarshal file contents: %v\n", err)
		return err
	}

	return nil
}

func (s *FileStorage) Close() {
	// Save the current metrics to the file before closing the storage
	if err := s.SaveToFile(); err != nil {
		logger.Log.Logger.WithError(err).Error("Failed to save metrics to file")
		// Handle error
	}
}

func (s *FileStorage) SaveToFile() error {
	data, err := json.MarshalIndent(s.memStorage, "", "  ")
	if err != nil {
		return err
	}

	// fmt.Println("Saving metrics to file:", s.fileStoragePath)
	// fmt.Println("File contents:", string(data))

	err = os.WriteFile(s.fileStoragePath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) saveAtInterval() {
	ticker := time.NewTicker(time.Duration(s.storeInterval) * time.Second)
	defer ticker.Stop()

	quit := make(chan struct{})
	defer close(quit)

	for {
		select {
		case <-ticker.C:
			if err := s.SaveToFile(); err != nil {
				logger.Log.Logger.WithError(err).Error("Failed to save metrics to file")
			}
		case <-quit:
			return
		}
	}
}

// Implement StorageInterface methods by delegating to memStorage
func (s *FileStorage) SetGauge(key string, value float64) {
	s.memStorage.SetGauge(key, value)
}

func (s *FileStorage) GetGauge(key string) (float64, bool) {
	return s.memStorage.GetGauge(key)
}

func (s *FileStorage) SetCounter(key string, value int64) {
	s.memStorage.SetCounter(key, value)
}

func (s *FileStorage) GetCounter(key string) (int64, bool) {
	return s.memStorage.GetCounter(key)
}

func (s *FileStorage) GetAllGauges() map[string]float64 {
	return s.memStorage.GetAllGauges()
}

func (s *FileStorage) GetAllCounters() map[string]int64 {
	return s.memStorage.GetAllCounters()
}

func (s *FileStorage) GetAllMetrics() map[string]interface{} {
	return s.memStorage.GetAllMetrics()
}
