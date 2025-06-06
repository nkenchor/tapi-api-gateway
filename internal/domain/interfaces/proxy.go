// internal/domain/interfaces/proxy.go
package interfaces

import (
	"net/http"
	"tapi-api-gateway-service/internal/domain/models"
)

// ProxyService interface defines methods for proxying requests to KrakenD
type ProxyService interface {
	ProxyRequest(req *http.Request) (*http.Response, *models.ApplicationError)
}
