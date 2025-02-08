package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// LogEntry represents the structure of a log entry.
type LogEntry struct {
	Message    string `json:"message"`
	RequestID  string `json:"request_id"`
	UserID     string `json:"user_id"`
	IP         string `json:"ip_address"`
	StatusCode int    `json:"status_code"`
	Duration   int64  `json:"duration"`
	Error      string `json:"error"`
}

// InitLogger configures the logger for structured JSON logging.
func InitLogger() {
	log = logrus.New()
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Create log file with timestamp
	currentTime := time.Now().Format("2006-01-02")
	file, err := os.OpenFile(
		filepath.Join("logs", currentTime+".log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.TraceLevel) // Default level can be adjusted
	//log.SetReportCaller(true)       // Include caller information (filename, function, line number)
}

// GenerateStackTrace generates a dynamic stack trace
func GenerateStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// getCallerInfo retrieves the caller's function name, filename, and line number.
// func getCallerInfo() (string, string, int) {
// 	pc, filename, line, ok := runtime.Caller(2)
// 	if !ok {
// 		return "", "", 0
// 	}
// 	funcName := runtime.FuncForPC(pc).Name()
// 	return funcName, filename, line
// }

// logMessage is a generic function to log messages at different levels.
func logMessage(level logrus.Level, message string, logFields logrus.Fields, errorMessage string) {
	// Add error and stack trace info if it's an error log
	if level == logrus.ErrorLevel || level == logrus.FatalLevel || level == logrus.PanicLevel {
		logFields["error"] = errorMessage
		logFields["stack_trace"] = GenerateStackTrace()
	}
	log.WithFields(logFields).Log(level, message)
}

// LogRequest logs the HTTP request details at INFO level.
func LogRequest(requestID, userID, ip string, statusCode int, duration int64, errorMessage string) {
	logEntry := LogEntry{
		Message:    "Request processed",
		RequestID:  requestID,
		UserID:     userID,
		IP:         ip,
		StatusCode: statusCode,
		Duration:   duration,
		Error:      errorMessage,
	}

	logMessage(logrus.TraceLevel, logEntry.Message, logrus.Fields{
		"request_id":  logEntry.RequestID,
		"user_id":     logEntry.UserID,
		"ip_address":  logEntry.IP,
		"status_code": logEntry.StatusCode,
		"duration":    logEntry.Duration,
		"error":       logEntry.Error,
	}, logEntry.Error)
}

// LogInfo logs general info messages at INFO level.
func LogInfo(message string) {
	logMessage(logrus.InfoLevel, message, logrus.Fields{}, "")
}

// LogWarn logs warnings at WARN level.
func LogWarn(message string) {
	logMessage(logrus.WarnLevel, message, logrus.Fields{}, "")
}

// LogDebug logs debug messages at DEBUG level.
func LogDebug(message string) {
	logMessage(logrus.DebugLevel, message, logrus.Fields{}, "")
}

// LogFatal logs fatal messages at FATAL level.
func LogFatal(message string) {
	logMessage(logrus.FatalLevel, message, logrus.Fields{}, "")
}

// LogError logs error messages at FATAL level.
func LogError(message string) {
	logMessage(logrus.ErrorLevel, message, logrus.Fields{}, "")
}

// LogError logs error messages at FATAL level.
func LogPanic(message string) {
	logMessage(logrus.PanicLevel, message, logrus.Fields{}, "")
}
