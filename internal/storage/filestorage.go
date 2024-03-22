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

func NewFileStorage(memStorage *MemStorage, storeInterval int, fileStoragePath string, restore bool) *FileStorage {
	storage := &FileStorage{
		memStorage:      memStorage,
		storeInterval:   storeInterval,
		fileStoragePath: fileStoragePath,
		restore:         restore,
	}

	// Load previously saved metrics from the file at startup
	if storage.restore {
		if err := storage.LoadFromFile(); err != nil {
			logger.Log.Logger.WithError(err).Error("Failed to load metrics from file")
			// log.Logger.WithError(err).Error("Failed to load metrics from file")
		}
	}

	// Save current metrics to the file at the specified interval
	go storage.saveAtInterval()

	return storage
}

func (s *FileStorage) LoadFromFile() error {
	file, err := os.ReadFile(s.fileStoragePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, s.memStorage)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) SaveToFile() error {
	data, err := json.Marshal(s.memStorage)
	if err != nil {
		return err
	}

	fmt.Println("Saving metrics to file:", s.fileStoragePath)

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
