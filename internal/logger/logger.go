package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func New(logFile string) error {
	Log = logrus.New()

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}

	Log.Out = file
	Log.SetFormatter(&logrus.JSONFormatter{})

	return nil
}

func Close() error {
	if Log != nil && Log.Out != nil {
		if file, ok := Log.Out.(*os.File); ok {
			return file.Close()
		}
	}
	return nil
}
