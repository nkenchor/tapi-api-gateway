// internal/infrastructure/config/config_provider.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"tapi-api-gateway-service/internal/domain/models"

	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
)

// ConfigProvider handles fetching configuration from environment variables and Consul
type ConfigProvider struct {
	consulClient *api.Client
}

// NewConfigProvider initializes a new ConfigProvider
func NewConfigProvider() (*ConfigProvider, error) {
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
		return nil, fmt.Errorf("failed to create Consul client: %v", err)
	}

	return &ConfigProvider{consulClient: client}, nil
}

// LoadConfig retrieves configuration values
func (p *ConfigProvider) LoadConfig() (*models.Config, error) {
	config := &models.Config{
		Port:                  getEnvAsInt("PORT", 30300),
		KrakenDURL:            p.getConsulValue("KRAKEND_URL", "http://localhost:30000"),
		AppName:               getEnv("APP_NAME", "tapi-api-gateway-service"),
		LogFile:               fmt.Sprintf("%s.log", getEnv("APP_NAME", "tapi-api-gateway-service")),
		LogDir:                p.getConsulValue("LOG_DIR", "logs"),
		LaunchURL:             p.getConsulValue("LAUNCH_URL", "swagger"),
		Address:               p.getConsulValue("ADDRESS", "localhost"),
		Mode:                  p.getConsulValue("MODE", "production"),
		DBType:                p.getConsulValue("DB_TYPE", "mongodb"),
		DBURL:                 p.getConsulValue("DB_URL", "mongodb://localhost:27017/tapi"),
		DB:                    p.getConsulValue("DB", "tapi"),
		RedisHost:             p.getConsulValue("REDIS_HOST", "127.0.0.1"),
		RedisPort:             p.getConsulValue("REDIS_PORT", "6379"),
		VaultURL:              p.getConsulValue("VAULT_URL", "http://localhost:8200"),
		VaultToken:            p.getConsulValue("VAULT_TOKEN", ""),
		VaultPath:             p.getConsulValue("VAULT_PATH", "tapi"),
		VaultKey:              p.getConsulValue("VAULT_KEY", ""),
		SecretKey:             p.getConsulValue("SECRET_KEY", ""),
		JWTKey:                p.getConsulValue("JWT_KEY", ""),
		JWTAudience:           p.getConsulValue("JWT_AUDIENCE", "tapi"),
		JWTIssuer:             p.getConsulValue("JWT_ISSUER", "https://localhost:3003"),
		JWTExpiry:             p.getConsulValueAsInt("JWT_EXPIRY", 600),
		OTPExpiry:             p.getConsulValueAsInt("OTP_EXPIRY", 60),
		OTPDigits:             p.getConsulValueAsInt("OTP_DIGITS", 6),
		TTL:                   p.getConsulValueAsInt("TTL", 900),
		KeycloakURL:           p.getConsulValue("KEYCLOAK_URL", "http://localhost:8181"),
		KeycloakAdminUsername: p.getConsulValue("KEYCLOAK_ADMIN_USERNAME", "admin"),
		KeycloakAdminPassword: p.getConsulValue("KEYCLOAK_ADMIN_PASSWORD", "admin"),
		KeycloakRealm:         p.getConsulValue("KEYCLOAK_REALM", "tapi"),
		KeycloakClientName:    p.getConsulValue("KEYCLOAK_CLIENT_NAME", "tapi-identity-service"),
		KeycloakClientID:      p.getConsulValue("KEYCLOAK_CLIENT_ID", ""),
		KeycloakClientSecret:  p.getConsulValue("KEYCLOAK_CLIENT_SECRET", ""),
		SMTPUsername:          p.getConsulValue("SMTP_USERNAME", "notification@tapi.com"),
		SMTPPassword:          p.getConsulValue("SMTP_PASSWORD", "notification"),
		SMTPHost:              p.getConsulValue("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:              p.getConsulValueAsInt("SMTP_PORT", 587),
	}
	return config, nil
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

func (p *ConfigProvider) getConsulValue(key string, defaultValue string) string {
	kvPair, _, err := p.consulClient.KV().Get(fmt.Sprintf("config/%s", key), nil)
	if err != nil || kvPair == nil || len(kvPair.Value) == 0 {
		fmt.Printf("Error fetching %s from Consul: %v. Using default value.\n", key, err)
		return defaultValue
	}
	return string(kvPair.Value)
}

func (p *ConfigProvider) getConsulValueAsInt(key string, defaultValue int) int {
	valueStr := p.getConsulValue(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
