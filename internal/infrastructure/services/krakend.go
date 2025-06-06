package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"
	
)


var krakendCmd *exec.Cmd


func StartKrakenD(loggingService interfaces.LoggingService) error {
	// Kill any process running on port 30000
	if err := KillProcessOnPort(30000, loggingService); err != nil {
		return err
	}

	// Check if KrakenD is already running
	if IsKrakenDRunning() {
		loggingService.Log(models.LogData{
			Level:   "INFO",
			Message: "KrakenD is already running, stopping existing process...",
			AppName: "tapi-api-gateway-service",
		})

		// Stop the currently running KrakenD process
		if err := StopKrakenD(loggingService); err != nil {
			loggingService.LogError(models.LogData{
				Level:   "ERROR",
				Message: "Failed to stop existing KrakenD process before starting a new one",
				AppName: "tapi-api-gateway-service",
			}, err)
			return err
		}
	}

	// Start a new KrakenD process
	fmt.Println("Starting KrakenD...")
	cmd := exec.Command("krakend", "run", "-c", "configuration/krakend.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		loggingService.LogError(models.LogData{
			Level:   "ERROR",
			Message: "Failed to start KrakenD process",
			AppName: "tapi-api-gateway-service",
		}, err)
		return err
	}
	krakendCmd = cmd // Save the running command to the global variable
	loggingService.Log(models.LogData{
		Level:   "INFO",
		Message: fmt.Sprintf("KrakenD started with PID %d", cmd.Process.Pid),
		AppName: "tapi-api-gateway-service",
	})
	return nil
}

func StopKrakenD(loggingService interfaces.LoggingService) error {
	if krakendCmd != nil && krakendCmd.Process != nil {
		if err := krakendCmd.Process.Kill(); err != nil {
			loggingService.LogError(models.LogData{
				Level:   "ERROR",
				Message: "Failed to kill KrakenD process",
				AppName: "tapi-api-gateway-service",
			}, err)
			return err
		}
	}
	return nil
}

func IsKrakenDRunning() bool {
	// Check if the krakendCmd process is still alive
	if krakendCmd != nil && krakendCmd.Process != nil {
		err := krakendCmd.Process.Signal(syscall.Signal(0))
		return err == nil
	}
	return false
}

func WaitForKrakenD() error {
	// Wait up to 10 seconds for KrakenD to start
	for i := 0; i < 5; i++ {
		if IsKrakenDRunning() {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("KrakenD did not start within the expected time")
}

func GenerateKrakenDConfig(loggingService interfaces.LoggingService) error {
	templatePath := "configuration/template/api_gateway_template.json"
	servicesPath := "configuration/services"
	outputPath := "configuration/krakend.json"

	// Load the template configuration
	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		loggingService.LogError(models.LogData{
			Level:   "ERROR",
			Message: "Failed to read template file",
			AppName: "tapi-api-gateway-service",
		}, err)
		return err
	}

	var templateConfig map[string]interface{}
	if err := json.Unmarshal(templateData, &templateConfig); err != nil {
		loggingService.LogError(models.LogData{
			Level:   "ERROR",
			Message: "Failed to parse template JSON",
			AppName: "tapi-api-gateway-service",
		}, err)
		return err
	}

	// Initialize the endpoints array
	endpoints := make([]interface{}, 0)

	// Load all service configuration files
	serviceFiles, err := filepath.Glob(filepath.Join(servicesPath, "*.json"))
	if err != nil {
		loggingService.LogError(models.LogData{
			Level:   "ERROR",
			Message: "Failed to read service configuration files",
			AppName: "tapi-api-gateway-service",
		}, err)
		return err
	}

	// Combine service configurations into the endpoints
	for _, serviceFile := range serviceFiles {
		serviceData, err := ioutil.ReadFile(serviceFile)
		if err != nil {
			loggingService.LogError(models.LogData{
				Level:   "ERROR",
				Message: fmt.Sprintf("Failed to read service configuration file: %s", serviceFile),
				AppName: "tapi-api-gateway-service",
			}, err)
			return err
		}

		var serviceEndpoints []interface{}
		if err := json.Unmarshal(serviceData, &serviceEndpoints); err != nil {
			loggingService.LogError(models.LogData{
				Level:   "ERROR",
				Message: fmt.Sprintf("Failed to parse service configuration JSON: %s", serviceFile),
				AppName: "tapi-api-gateway-service",
			}, err)
			return err
		}

		// Append service endpoints to the endpoints array
		endpoints = append(endpoints, serviceEndpoints...)
	}

	// Update the template's endpoints with the combined endpoints
	templateConfig["endpoints"] = endpoints

	// Write the merged configuration to the output file
	outputData, err := json.MarshalIndent(templateConfig, "", "  ")
	if err != nil {
		loggingService.LogError(models.LogData{
			Level:   "ERROR",
			Message: "Failed to marshal combined JSON",
			AppName: "tapi-api-gateway-service",
		}, err)
		return err
	}

	if err := ioutil.WriteFile(outputPath, outputData, 0644); err != nil {
		loggingService.LogError(models.LogData{
			Level:   "ERROR",
			Message: "Failed to write combined JSON to output file",
			AppName: "tapi-api-gateway-service",
		}, err)
		return err
	}

	loggingService.Log(models.LogData{
		Level:   "INFO",
		Message: "Successfully generated krakend.json",
		AppName: "tapi-api-gateway-service",
	})

	return nil
}
