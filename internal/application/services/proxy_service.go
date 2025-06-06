// internal/application/services/proxy_service.go
package services

import (
	"net/http"
	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"
)

// ProxyServiceImpl is the implementation of the ProxyService interface
type ProxyServiceImpl struct{}

// NewProxyService creates a new instance of ProxyServiceImpl
func NewProxyService() interfaces.ProxyService {
	return &ProxyServiceImpl{}
}

// ProxyRequest forwards a request to KrakenD and returns the response
func (s *ProxyServiceImpl) ProxyRequest(req *http.Request) (*http.Response, *models.ApplicationError) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, models.NewApplicationError(models.ProxyError, "Failed to proxy request", http.StatusInternalServerError, nil)
	}
	return resp, nil
}
