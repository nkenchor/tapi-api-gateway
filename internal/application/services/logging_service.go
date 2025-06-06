// internal/application/services/logging_service.go
package services

import (
	"fmt"
	"os"
	"path/filepath"
	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"

	"github.com/sirupsen/logrus"
)

// LoggingServiceImpl is the implementation of the LoggingService interface
type LoggingServiceImpl struct {
	logger *logrus.Logger
}

// NewLoggingService creates a new instance of LoggingServiceImpl
func NewLoggingService() interfaces.LoggingService {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Ensure the logs directory exists
	logDir := "logs"
	logFile := filepath.Join(logDir, "tapi-api-gateway-service.log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
	} else {
		// Create or open the log file
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open log file: %v\n", err)
		} else {
			// Set log output to the file
			logger.SetOutput(file)
			logger.SetFormatter(&logrus.JSONFormatter{})
			logger.SetLevel(logrus.InfoLevel)
		}
	}

	return &LoggingServiceImpl{logger: logger}
}

// Log logs the given data
func (s *LoggingServiceImpl) Log(data models.LogData) error {
	if s == nil || s.logger == nil {
		return fmt.Errorf("logging service or logger is not initialized")
	}

	logFields := logrus.Fields{
		"timestamp":        data.Timestamp,
		"level":            data.Level,
		"method":           data.Method,
		"status_code":      data.StatusCode,
		"path":             data.Path,
		"forwarded_for":    data.ForwardedFor,
		"response_time":    data.ResponseTime,
		"version":          data.Version,
		"correlation_id":   data.CorrelationID,
		"app_name":         data.AppName,
		"application_host": data.ApplicationHost,
		"ip":               data.IP,
		"computerName":     data.ComputerName,
		"logId":            data.LogID,
		"message":          data.Message,
	}

	s.logger.WithFields(logFields).Info(data.Message)
	return nil
}

// LogError logs an error
func (s *LoggingServiceImpl) LogError(data models.LogData, err error) error {
	if s == nil || s.logger == nil {
		return fmt.Errorf("logging service or logger is not initialized")
	}

	logFields := logrus.Fields{
		"timestamp":        data.Timestamp,
		"level":            "ERROR",
		"method":           data.Method,
		"status_code":      data.StatusCode,
		"path":             data.Path,
		"forwarded_for":    data.ForwardedFor,
		"response_time":    data.ResponseTime,
		"version":          data.Version,
		"correlation_id":   data.CorrelationID,
		"app_name":         data.AppName,
		"application_host": data.ApplicationHost,
		"ip":               data.IP,
		"computerName":     data.ComputerName,
		"logId":            data.LogID,
		"message":          data.Message,
	}

	// Ensure the error is captured in the log if it's not nil
	if err != nil {
		logFields["error"] = fmt.Sprintf("%v", err)
	}

	s.logger.WithFields(logFields).Error(data.Message)
	return nil
}
