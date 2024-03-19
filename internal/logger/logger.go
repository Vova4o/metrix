package logger

import (
    "fmt"
    "os"

    "github.com/sirupsen/logrus"
)

type LoggerInterface interface {
    CreateLogger(name string) error
    CloseLogger() error
    SetOutput()
}

func NewFileLogger(name string) (*FileLogger, error) {
    var logger FileLogger
    err := logger.CreateLogger(name)
    if err != nil {
        return nil, err
    }
    return &logger, nil
}

type FileLogger struct {
    LogFile *os.File
    Logger  *logrus.Logger
}

func (l *FileLogger) CreateLogger(name string) error {
    var err error
    l.LogFile, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
    if err != nil {
        fmt.Printf("Failed to open log file: %v\n", err)
        return err
    }

    l.Logger = logrus.New()
    l.Logger.Out = l.LogFile
    l.Logger.Formatter = &logrus.JSONFormatter{}

    return nil
}

func (l *FileLogger) CloseLogger() error {
    if l.LogFile != nil {
        err := l.LogFile.Close()
        if err != nil {
            return err
        }
    }
    return nil
}

func (l *FileLogger) SetOutput() {
    logrus.SetOutput(l.LogFile)
}