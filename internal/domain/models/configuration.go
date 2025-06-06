// internal/domain/models/config.go
package models

// Config holds all the configuration values for the application
type Config struct {
	Port                  int
	KrakenDURL            string
	AppName               string
	LogFile               string
	LogDir                string
	LaunchURL             string
	Address               string
	Mode                  string
	DBType                string
	DBURL                 string
	DB                    string
	RedisHost             string
	RedisPort             string
	VaultURL              string
	VaultToken            string
	VaultPath             string
	VaultKey              string
	SecretKey             string
	JWTKey                string
	JWTAudience           string
	JWTIssuer             string
	JWTExpiry             int
	OTPExpiry             int
	OTPDigits             int
	TTL                   int
	KeycloakURL           string
	KeycloakAdminUsername string
	KeycloakAdminPassword string
	KeycloakRealm         string
	KeycloakClientName    string
	KeycloakClientID      string
	KeycloakClientSecret  string
	SMTPUsername          string
	SMTPPassword          string
	SMTPHost              string
	SMTPPort              int
}
