// internal/domain/models/log_data.go
package models

import (
	"time"
)

// LogData represents the structured logging data for each request
type LogData struct {
	Timestamp       string            `json:"timestamp"`
	Level           string            `json:"level"`
	Method          string            `json:"method"`
	StatusCode      string            `json:"status_code"`
	Path            string            `json:"path"`
	ForwardedFor    string            `json:"forwarded_for"`
	ResponseTime    string            `json:"response_time"`
	Version         string            `json:"version"`
	CorrelationID   string            `json:"correlation_id"`
	AppName         string            `json:"app_name"`
	ApplicationHost string            `json:"application_host"`
	IP              string            `json:"ip"`
	ComputerName    string            `json:"computerName"`
	LogID           string            `json:"logId"`
	Message         string            `json:"message"`
	RequestHeaders  map[string]string `json:"requestHeaders"`
	RequestQuery    map[string]string `json:"requestQuery"`
	RequestBody     string            `json:"requestBody"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	ResponseBody    string            `json:"responseBody"`
	RequestDuration float64           `json:"requestDuration"`
	Error           string            `json:"error"`
	Stack           string            `json:"stack"`
	StartTime       time.Time         `json:"-"`
}
