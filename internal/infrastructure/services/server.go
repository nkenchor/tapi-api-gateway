package services

import (
	"fmt"
	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"
	"tapi-api-gateway-service/internal/infrastructure/controllers"
	"tapi-api-gateway-service/internal/infrastructure/proxy"

	"github.com/gin-gonic/gin"
)

func StartGinServer(port int, loggingService interfaces.LoggingService, serverUpSignalChan chan struct{}) {
	// Initialize proxy and health services
	proxyService := proxy.NewProxyService()
	proxyController := controllers.NewProxyController(proxyService, loggingService)

	// Create a new Gin router
	r := gin.Default()

	// Middleware to log requests
	r.Use(gin.Logger(), gin.Recovery())

	// Define the catch-all route for proxy requests
	r.Any("/*proxyPath", proxyController.ProxyRequest)

	// Start the Gin server
	go func() {
		loggingService.Log(models.LogData{
			Level:   "INFO",
			Message: fmt.Sprintf("Proxy server running on port %d", port),
			AppName: "tapi-api-gateway-service",
		})
		if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
			loggingService.LogError(models.LogData{
				Level:   "ERROR",
				Message: "Failed to run Gin proxy",
				AppName: "tapi-api-gateway-service",
			}, err)
		
		}
	}()

	// Notify that the server is running
	serverUpSignalChan <- struct{}{}
}
