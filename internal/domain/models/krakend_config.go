// internal/domain/models/krakend_config.go
package models

// KrakenDConfig represents the structure of the KrakenD configuration
type KrakenDConfig struct {
	Version   int                      `json:"version"`
	Name      string                   `json:"name"`
	Port      int                      `json:"port"`
	Endpoints []map[string]interface{} `json:"endpoints"`
}
