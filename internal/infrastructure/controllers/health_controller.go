// internal/infrastructure/controllers/health_controller.go
package controllers

import (
	"net/http"

	"tapi-api-gateway-service/internal/domain/interfaces"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	healthService interfaces.HealthService
}

// NewHealthController creates a new instance of HealthController
func NewHealthController(healthService interfaces.HealthService) *HealthController {
	return &HealthController{healthService: healthService}
}

// CheckHealth handles the health check endpoint
func (hc *HealthController) CheckHealth(c *gin.Context) {
	// Call the health service to perform a health check
	status, err := hc.healthService.CheckHealth()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "message": status})
}
