package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"Vova4o/metrix/internal/handlers"
	"Vova4o/metrix/internal/logger"
)

type FileStorage struct {
	handlers.Storager
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func NewFile(s handlers.Storager, storeInterval int, fileStoragePath string, restore bool) (*FileStorage, error) {
	fs := &FileStorage{
		Storager:        s,
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
			logger.Log.WithError(err).Error("Failed to load metrics from file")
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
	_, err := os.Stat(s.fileStoragePath)
	if os.IsNotExist(err) {
		file, err := os.Create(s.fileStoragePath)
		if err != nil {
			return fmt.Errorf("failed to create file at path %s: %w", s.fileStoragePath, err)
		}
		defer file.Close()
	} else if err != nil {
		return fmt.Errorf("failed to check if file exists at path %s: %w", s.fileStoragePath, err)
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

	err = json.Unmarshal(contents, s.Storager)
	if err != nil {
		fmt.Printf("Failed to unmarshal file contents: %v\n", err)
		return err
	}

	return nil
}

func (s *FileStorage) Close() {
	// Save the current metrics to the file before closing the storage
	if err := s.SaveToFile(); err != nil {
		logger.Log.WithError(err).Error("Failed to save metrics to file")
		// Handle error
	}
}

func (s *FileStorage) SaveToFile() error {
	data, err := json.MarshalIndent(s.Storager, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println("Saving metrics to file:", s.fileStoragePath)
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

	for range ticker.C {
		if err := s.SaveToFile(); err != nil {
			logger.Log.WithError(err).Error("Failed to save metrics to file")
		}
	}
}
