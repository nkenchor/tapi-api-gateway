package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"
)

// ProxyServiceImpl implements the ProxyService interface
type ProxyServiceImpl struct{}

// NewProxyService creates a new instance of ProxyServiceImpl
func NewProxyService() interfaces.ProxyService {
	return &ProxyServiceImpl{}
}

// ProxyRequest forwards the HTTP request to KrakenD and handles the response
func (s *ProxyServiceImpl) ProxyRequest(req *http.Request) (*http.Response, *models.ApplicationError) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Return an ApplicationError for any network-related error
		return nil, models.NewApplicationError(
			models.ProxyError, 
			"Failed to proxy request",
			http.StatusInternalServerError,
			nil,
		)
	}

	// Check if the response status code indicates an error
	if resp.StatusCode >= 400 {
		// Read the response body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, models.NewApplicationError(
				models.ProxyError,
				"Failed to read response body from KrakenD",
				http.StatusInternalServerError,
				nil,
			)
		}

		// Restore the response body to allow reading it later
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Try to parse the response body as JSON
		var jsonBody map[string]interface{}
		if jsonErr := json.Unmarshal(bodyBytes, &jsonBody); jsonErr != nil {
			// If the response is not valid JSON, include the raw response body as a string
			if len(bodyBytes) == 0 {
				return nil, models.StatusCodeToApplicationError(resp.StatusCode)
			}
			return nil, models.NewApplicationError(
				models.ProxyError,
				fmt.Sprintf("KrakenD returned an error response with status code %d: %s", resp.StatusCode, string(bodyBytes)),
				resp.StatusCode,
				nil,
			)
		}

	}

	return resp, nil
}
