// internal/infrastructure/logging/log_provider.go
package logging

import (
	"os"
	"tapi-api-gateway-service/internal/domain/models"

	"github.com/sirupsen/logrus"
)

// LogProvider handles logging using Logrus
type LogProvider struct {
	logger *logrus.Logger
}

// NewLogProvider initializes a new LogProvider
func NewLogProvider(logFile string, logDir string) (*LogProvider, error) {
	// Ensure log directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}
	}

	// Set up log file
	logFilePath := logDir + "/" + logFile
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Configure Logrus
	logger := logrus.New()
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return &LogProvider{logger: logger}, nil
}

// Log writes a log entry
func (p *LogProvider) Log(data models.LogData) {
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

	p.logger.WithFields(logFields).Info(data.Message)
}

// LogError writes an error log entry
func (p *LogProvider) LogError(data models.LogData, err error) {
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
		"error":            err.Error(),
	}

	p.logger.WithFields(logFields).Error(data.Message)
}
