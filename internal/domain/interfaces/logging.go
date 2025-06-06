// internal/domain/interfaces/logging.go
package interfaces

import "tapi-api-gateway-service/internal/domain/models"

// LoggingService interface defines methods for structured logging
type LoggingService interface {
	Log(data models.LogData) error
	LogError(data models.LogData, err error) error
}
