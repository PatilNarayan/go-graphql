// logger/logger.go
package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		Log.Fatal("Failed to create logs directory:", err)
	}

	// Create log file with timestamp
	currentTime := time.Now().Format("2006-01-02")
	file, err := os.OpenFile(
		filepath.Join("logs", currentTime+".log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		Log.Fatal("Failed to open log file:", err)
	}

	Log.SetOutput(file)
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}

func AddContext(err error) *logrus.Entry {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		return Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"file":  file,
			"line":  line,
		})
	}
	return Log.WithFields(logrus.Fields{})
}
