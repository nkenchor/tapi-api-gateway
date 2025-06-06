// internal/application/services/config_service.go
package services

import (
	"fmt"
	"os"
	"strconv"
	"tapi-api-gateway-service/internal/domain/interfaces"
	"tapi-api-gateway-service/internal/domain/models"

	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
)

// ConfigServiceImpl is the implementation of the ConfigService interface
type ConfigServiceImpl struct {
	consulClient *api.Client
}

// NewConfigService creates a new instance of ConfigServiceImpl
func NewConfigService() (interfaces.ConfigService, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found or could not be loaded.")
	}

	// Initialize Consul client
	consulConfig := api.DefaultConfig()
	consulToken := os.Getenv("CONSUL_HTTP_TOKEN")
	if consulToken != "" {
		consulConfig.Token = consulToken
	} else {
		fmt.Println("Warning: CONSUL_HTTP_TOKEN environment variable is not set. Attempting to use unauthenticated Consul client.")
	}

	client, err := api.NewClient(consulConfig)
	if err != nil {
		// If Consul client creation fails, use a nil client and print a warning
		fmt.Println("Failed to create Consul client. Falling back to environment variables and default values.")
		client = nil // Marking client as nil to handle fallback in GetConfig
	}

	return &ConfigServiceImpl{consulClient: client}, nil
}

// GetConfig retrieves configuration values
func (s *ConfigServiceImpl) GetConfig() (*models.Config, error) {
	// If consulClient is nil, handle the fallback
	if s.consulClient == nil {
		fmt.Println("Using default values and environment variables due to missing Consul client.")
	}

	config := &models.Config{
		Port:                  getEnvAsInt("PORT", 30400),
		KrakenDURL:            s.getConsulOrEnvValue("KRAKEND_URL", "http://localhost:30000"),
		AppName:               getEnv("APP_NAME", "tapi-api-gateway-service"),
		LogFile:               fmt.Sprintf("%s.log", getEnv("APP_NAME", "tapi-api-gateway-service")),
		LogDir:                s.getConsulOrEnvValue("LOG_DIR", "logs"),
		LaunchURL:             s.getConsulOrEnvValue("LAUNCH_URL", "swagger"),
		Address:               s.getConsulOrEnvValue("ADDRESS", "localhost"),
		Mode:                  s.getConsulOrEnvValue("MODE", "production"),
		DBType:                s.getConsulOrEnvValue("DB_TYPE", "mongodb"),
		DBURL:                 s.getConsulOrEnvValue("DB_URL", "mongodb://localhost:27017/tapi"),
		DB:                    s.getConsulOrEnvValue("DB", "tapi"),
		RedisHost:             s.getConsulOrEnvValue("REDIS_HOST", "127.0.0.1"),
		RedisPort:             s.getConsulOrEnvValue("REDIS_PORT", "6379"),
		VaultURL:              s.getConsulOrEnvValue("VAULT_URL", "http://localhost:8200"),
		VaultToken:            s.getConsulOrEnvValue("VAULT_TOKEN", ""),
		VaultPath:             s.getConsulOrEnvValue("VAULT_PATH", "tapi"),
		VaultKey:              s.getConsulOrEnvValue("VAULT_KEY", ""),
		SecretKey:             s.getConsulOrEnvValue("SECRET_KEY", ""),
		JWTKey:                s.getConsulOrEnvValue("JWT_KEY", ""),
		JWTAudience:           s.getConsulOrEnvValue("JWT_AUDIENCE", "tapi"),
		JWTIssuer:             s.getConsulOrEnvValue("JWT_ISSUER", "https://localhost:3003"),
		JWTExpiry:             s.getConsulOrEnvValueAsInt("JWT_EXPIRY", 600),
		OTPExpiry:             s.getConsulOrEnvValueAsInt("OTP_EXPIRY", 60),
		OTPDigits:             s.getConsulOrEnvValueAsInt("OTP_DIGITS", 6),
		TTL:                   s.getConsulOrEnvValueAsInt("TTL", 900),
		KeycloakURL:           s.getConsulOrEnvValue("KEYCLOAK_URL", "http://localhost:8181"),
		KeycloakAdminUsername: s.getConsulOrEnvValue("KEYCLOAK_ADMIN_USERNAME", "admin"),
		KeycloakAdminPassword: s.getConsulOrEnvValue("KEYCLOAK_ADMIN_PASSWORD", "admin"),
		KeycloakRealm:         s.getConsulOrEnvValue("KEYCLOAK_REALM", "tapi"),
		KeycloakClientName:    s.getConsulOrEnvValue("KEYCLOAK_CLIENT_NAME", "tapi-identity-service"),
		KeycloakClientID:      s.getConsulOrEnvValue("KEYCLOAK_CLIENT_ID", ""),
		KeycloakClientSecret:  s.getConsulOrEnvValue("KEYCLOAK_CLIENT_SECRET", ""),
		SMTPUsername:          s.getConsulOrEnvValue("SMTP_USERNAME", "notification@tapi.com"),
		SMTPPassword:          s.getConsulOrEnvValue("SMTP_PASSWORD", "notification"),
		SMTPHost:              s.getConsulOrEnvValue("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:              s.getConsulOrEnvValueAsInt("SMTP_PORT", 587),
	}
	return config, nil
}

// Helper function to get values from Consul or fall back to env or default
func (s *ConfigServiceImpl) getConsulOrEnvValue(key string, defaultValue string) string {
	if s.consulClient != nil {
		value := getConsulValue(s.consulClient, key, "")
		if value != "" {
			return value
		}
	}
	return getEnv(key, defaultValue)
}

func (s *ConfigServiceImpl) getConsulOrEnvValueAsInt(key string, defaultValue int) int {
	if s.consulClient != nil {
		value := getConsulValueAsInt(s.consulClient, key, defaultValue)
		if value != defaultValue {
			return value
		}
	}
	return getEnvAsInt(key, defaultValue)
}

// Helper functions to get values
func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getConsulValue(client *api.Client, key string, defaultValue string) string {
	kvPair, _, err := client.KV().Get(fmt.Sprintf("config/%s", key), nil)
	if err != nil || kvPair == nil || len(kvPair.Value) == 0 {
		fmt.Printf("Error fetching %s from Consul: %v. Using default value.\n", key, err)
		return defaultValue
	}
	return string(kvPair.Value)
}

func getConsulValueAsInt(client *api.Client, key string, defaultValue int) int {
	valueStr := getConsulValue(client, key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
