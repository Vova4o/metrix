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

func NewFileStorage(memStorager Storager, storeInterval int, fileStoragePath string, restore bool) (*FileStorage, error) {
	memStorage, ok := memStorager.(*MemStorage)
	if !ok {
		err := fmt.Errorf("expected *storage.MemStorage type")
		logger.Log.Logger.Errorf(err.Error())
		return nil, err
	}
	
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
        if err := fs.createFileIfNotExists(); err != nil {
            return nil, err
        }

        // File exists, load it into memory
        fmt.Println("Loading metrics from file:", fs.fileStoragePath)
        if err := fs.LoadFromFile(); err != nil {
            logger.Log.Logger.WithError(err).Error("Failed to load metrics from file")
            return nil, err
        }
    }

	// Save current metrics to the file at the specified interval
	if fs.fileStoragePath != "" {
		go fs.saveAtInterval()
	}

	return fs, nil
}

func (s *FileStorage) createFileIfNotExists() error {
    if _, err := os.Stat(s.fileStoragePath); os.IsNotExist(err) {
        // create a new file if it doesn't exist
        file, err := os.Create(s.fileStoragePath)
        if err != nil {
            logger.Log.Logger.WithError(err).Error("Failed to create new file")
            return err
        }
        defer file.Close()
        logger.Log.Logger.Info("No previous metrics file found. Created a new one.")
    } else if err != nil {
        // some other error occurred when trying to stat the file
        logger.Log.Logger.WithError(err).Error("Failed to check if metrics file exists")
        return err
    }
    return nil
}


func (s *FileStorage) LoadFromFile() error {
    // Open the file
    file, err := os.Open(s.fileStoragePath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Check if the file is empty
    info, err := file.Stat()
    if err != nil {
        return err
    }
    if info.Size() == 0 {
        // The file is empty, return without error
        return nil
    }

    contents, err := os.ReadFile(s.fileStoragePath)
    if err != nil {
        return err
    }

    // Check if the contents are empty
    if len(contents) == 0 {
        // The contents are empty, return without error
        return nil
    }

    err = json.Unmarshal(contents, s.memStorage)
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

    for {
        select {
        case <-ticker.C:
            if err := s.SaveToFile(); err != nil {
                logger.Log.Logger.WithError(err).Error("Failed to save metrics to file")
            }
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
