package controllers

import (
	"fmt"
	"io"
	"net/http"

	"strconv"

	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type ProxyController struct {
	proxyService   interfaces.ProxyService
	loggingService interfaces.LoggingService
}

// NewProxyController creates a new instance of ProxyController
func NewProxyController(proxyService interfaces.ProxyService, loggingService interfaces.LoggingService) *ProxyController {
	return &ProxyController{
		proxyService:   proxyService,
		loggingService: loggingService,
	}
}

// ProxyRequest handles forwarding of HTTP requests to KrakenD via the ProxyService
func (pc *ProxyController) ProxyRequest(c *gin.Context) {
	proxyPath := c.Param("proxyPath")

	// Ensure that proxyPath starts with a "/" and is not empty
	if proxyPath == "" {
		proxyPath = "/"
	} else if proxyPath[0] != '/' {
		proxyPath = "/" + proxyPath
	}

	// Add the query parameters from the original request to the target URL
	queryParams := c.Request.URL.RawQuery // Get the query string
	targetURL := "http://localhost:30000" + proxyPath
	if queryParams != "" {
		targetURL = targetURL + "?" + queryParams // Add query parameters if they exist
	}

	// Log the proxy path for debugging
	pc.loggingService.Log(models.LogData{
		Message: fmt.Sprintf("Proxying request to: %s", targetURL),
		Level:   "INFO",
		AppName: "tapi-api-gateway-service",
	})

	// Create a new request to be forwarded
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		pc.loggingService.LogError(models.LogData{
			Message: "Failed to create request",
			Level:   "ERROR",
			AppName: "tapi-api-gateway-service",
		}, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation error"})
		return
	}

	// Copy headers from the original request to the new request
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Forward the request using the proxy service
	resp, appErr := pc.proxyService.ProxyRequest(req)
	if appErr != nil {
		// Log the error and return a JSON response
		pc.loggingService.LogError(models.LogData{
			Message:    appErr.Error(),
			Level:      "ERROR",
			StatusCode: strconv.Itoa(appErr.StatusCode),
			AppName:    "tapi-api-gateway-service",
		}, nil)
		c.JSON(appErr.StatusCode, appErr) // Directly return the ApplicationError
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			pc.loggingService.LogError(models.LogData{
				Message: "Failed to close response body",
				Level:   "ERROR",
				AppName: "tapi-api-gateway-service",
			}, err)
		}
	}()

	// Copy the response headers to Gin context
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Set the status code
	c.Writer.WriteHeader(resp.StatusCode)

	// Copy the response body directly to Gin context
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		pc.loggingService.LogError(models.LogData{
			Message: "Failed to write response body",
			Level:   "ERROR",
			AppName: "tapi-api-gateway-service",
		}, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward response body"})
		return
	}

	// Log the successful response from KrakenD
	pc.loggingService.Log(models.LogData{
		Message:         fmt.Sprintf("Successfully proxied request to: %s", targetURL),
		Level:           "INFO",
		Method:          req.Method,
		Path:            req.URL.Path,
		StatusCode:      strconv.Itoa(resp.StatusCode),
		ForwardedFor:    c.ClientIP(),
		CorrelationID:   c.Request.Header.Get("X-Correlation-ID"),
		ApplicationHost: req.Host,
		AppName:         "tapi-api-gateway-service",
	})
}
