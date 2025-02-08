package gormlogger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// GORMLogger is a custom GORM logger that integrates with Logrus.
type GORMLogger struct {
	Logger *logrus.Logger
}

// NewGORMLogger creates a new instance of GORMLogger.
func NewGORMLogger() *GORMLogger {
	return &GORMLogger{
		Logger: logrus.New(), // Initialize Logrus Logger
	}
}

// LogMode sets the logging mode for GORM.
func (g *GORMLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

// Info logs info-level messages like SQL queries.
func (g *GORMLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	g.Logger.WithFields(logrus.Fields{
		"level":     "INFO",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   msg,
		"args":      args,
	}).Info("GORM Info: SQL query executed")
}

// Warn logs warnings (e.g., potential issues with queries).
func (g *GORMLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	g.Logger.WithFields(logrus.Fields{
		"level":     "WARN",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   msg,
		"args":      args,
	}).Warn("GORM Warning")
}

// Error logs errors, typically when queries fail.
func (g *GORMLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	g.Logger.WithFields(logrus.Fields{
		"level":     "ERROR",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   msg,
		"args":      args,
	}).Error("GORM Error")
}

// Trace logs the query duration, SQL query, and any errors.
func (g *GORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Calculate the query duration
	duration := time.Since(begin)

	// Get the SQL query and rows affected
	sql, rowsAffected := fc()

	// Log the query execution time and related details
	fields := logrus.Fields{
		"level":     "INFO",
		"timestamp": time.Now().Format(time.RFC3339),
		"duration":  duration.Milliseconds(),
		"sql":       sql,
		"rows":      rowsAffected,
	}

	if err != nil {
		// Log error if the query failed
		fields["error"] = err.Error()
		g.Logger.WithFields(fields).Error("GORM Query failed")
	} else {
		// Log successful query execution
		g.Logger.WithFields(fields).Info("GORM Query executed")
	}
}
