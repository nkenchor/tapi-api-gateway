// internal/application/services/health_service.go
package services

import (
	"fmt"
	"tapi-api-gateway-service/internal/domain/interfaces"
	"time"
)

// HealthServiceImpl is the implementation of the HealthService interface
type HealthServiceImpl struct{}

// NewHealthService creates a new instance of HealthServiceImpl
func NewHealthService() interfaces.HealthService {
	return &HealthServiceImpl{}
}

// CheckHealth checks the health status of the service
func (s *HealthServiceImpl) CheckHealth() (string, error) {
	// Add more comprehensive health checks as needed
	return fmt.Sprintf("The service is running like a Go Ninja! Uptime: %s", time.Since(startTime).String()), nil
}

var startTime = time.Now()
