// internal/domain/interfaces/service.go
package interfaces

import "tapi-api-gateway-service/internal/domain/models"

// ConfigService interface defines methods for fetching configurations
type ConfigService interface {
	GetConfig() (*models.Config, error)
}

// HealthService interface defines methods for checking health status
type HealthService interface {
	CheckHealth() (string, error)
}
