package services

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"

	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"
)

func KillProcessOnPort(port int, loggingService interfaces.LoggingService) error {
	loggingService.Log(models.LogData{
		Level:   "INFO",
		Message: fmt.Sprintf("Attempting to kill any process using port %d...", port),
		AppName: "tapi-api-gateway-service",
	})

	// Try to bind to the port
	address := fmt.Sprintf(":%d", port)
	conn, err := net.Listen("tcp", address)
	if err != nil {
		loggingService.Log(models.LogData{
			Level:   "INFO",
			Message: fmt.Sprintf("Port %d is already in use. Trying to kill the process using it...", port),
			AppName: "tapi-api-gateway-service",
		})

		// If the port is already in use, attempt to find and kill the process using the port
		return KillProcessOnPortByOS(port, loggingService)
	}

	// Close the connection if we successfully bind to the port
	defer conn.Close()

	loggingService.Log(models.LogData{
		Level:   "INFO",
		Message: fmt.Sprintf("Port %d is free to use.", port),
		AppName: "tapi-api-gateway-service",
	})

	return nil
}

// KillProcessOnPortByOS kills the process on a given port using the appropriate command for the operating system.
func KillProcessOnPortByOS(port int, loggingService interfaces.LoggingService) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		// Use lsof and kill on macOS and Linux
		cmd = exec.Command("bash", "-c", fmt.Sprintf("lsof -ti tcp:%d | xargs kill -9", port))
	case "windows":
		// Use netstat and taskkill on Windows
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("netstat -ano | findstr :%d && for /f \"tokens=5\" %%a in ('netstat -ano | findstr :%d') do taskkill /F /PID %%a", port, port))
	default:
		loggingService.Log(models.LogData{
			Level:   "ERROR",
			Message: fmt.Sprintf("Unsupported OS: %s", runtime.GOOS),
			AppName: "tapi-api-gateway-service",
		})
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// Run the command to kill the process
	if err := cmd.Run(); err != nil {
		loggingService.Log(models.LogData{
			Level:   "ERROR",
			Message: fmt.Sprintf("Failed to kill processes on port %d: %v", port, err),
			AppName: "tapi-api-gateway-service",
		})
		return err
	}

	loggingService.Log(models.LogData{
		Level:   "INFO",
		Message: fmt.Sprintf("Successfully killed process on port %d.", port),
		AppName: "tapi-api-gateway-service",
	})

	return nil
}
