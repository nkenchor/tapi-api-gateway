package main

import (

	"log"
	"os"
	"os/signal"
	"syscall"

	"tapi-api-gateway-service/internal/application/services"
	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"
	"tapi-api-gateway-service/internal/infrastructure/config"
	infrastructure "tapi-api-gateway-service/internal/infrastructure/services"
	
)

func main() {
	// Initialize logging service
	loggingService := services.NewLoggingService()
	loggingService.Log(models.LogData{
		Level:   "INFO",
		Message: "Starting tapi-api-gateway-service...",
		AppName: "tapi-api-gateway-service",
	})

	// Load configuration
	configProvider, err := config.NewConfigProvider()
	if err != nil {
		logFatal(loggingService, "Failed to initialize config provider", err)
	}
	appConfig, err := configProvider.LoadConfig()
	if err != nil {
		logFatal(loggingService, "Failed to load configuration", err)
	}

	// Generate the krakend.json configuration
	if err := infrastructure.GenerateKrakenDConfig(loggingService); err != nil {
		logFatal(loggingService, "Failed to generate KrakenD configuration", err)
	}

	// Channel to signal when the server is up
	serverUpSignalChan := make(chan struct{})

	// Start the Gin server in a goroutine
	go infrastructure.StartGinServer(appConfig.Port, loggingService, serverUpSignalChan)

	// Wait for the Gin server to start
	<-serverUpSignalChan

	// Ensure KrakenD is not already running and start it
	if infrastructure.IsKrakenDRunning() {
		infrastructure.StopKrakenD(loggingService)
	}
	if err := infrastructure.StartKrakenD(loggingService); err != nil {
		logFatal(loggingService, "Failed to start KrakenD", err)
	}

	// Wait for KrakenD to be fully operational
	if err := infrastructure.WaitForKrakenD(); err != nil {
		logFatal(loggingService, "KrakenD did not start properly", err)
	}
	loggingService.Log(models.LogData{
		Level:   "INFO",
		Message: "KrakenD started successfully.",
		AppName: "tapi-api-gateway-service",
	})

	// Handle graceful shutdown
	handleShutdown(loggingService)
}

func handleShutdown(loggingService interfaces.LoggingService) {
	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	<-signalChan

	loggingService.Log(models.LogData{
		Level:   "INFO",
		Message: "Received shutdown signal, stopping KrakenD and Gin server...",
		AppName: "tapi-api-gateway-service",
	})

	// Stop KrakenD
	if err := infrastructure.StopKrakenD(loggingService); err != nil {
		loggingService.LogError(models.LogData{
			Level:   "ERROR",
			Message: "Failed to stop KrakenD",
			AppName: "tapi-api-gateway-service",
		}, err)
	} else {
		loggingService.Log(models.LogData{
			Level:   "INFO",
			Message: "KrakenD stopped successfully.",
			AppName: "tapi-api-gateway-service",
		})
	}

	os.Exit(0)
}

func logFatal(loggingService interfaces.LoggingService, message string, err error) {
	loggingService.LogError(models.LogData{
		Level:   "ERROR",
		Message: message,
		AppName: "tapi-api-gateway-service",
	}, err)
	log.Fatalf("%s: %v", message, err)
}
