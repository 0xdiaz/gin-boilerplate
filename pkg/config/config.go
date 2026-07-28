package config

import (
	"fmt"
	"strings"

	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/spf13/viper"
)

// Environment constants
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

var appConfig *Configuration

// Get returns the loaded configuration. Call after SetupConfig().
// Returns nil if config has not been loaded.
func Get() *Configuration {
	return appConfig
}

// SetForTest sets the global configuration. For testing only; allows unit tests to run without .env.
func SetForTest(c *Configuration) {
	appConfig = c
}

// SetupConfig loads and validates configuration from .env file
// and stores it in the global config (accessible via Get()).
func SetupConfig() error {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error reading config file, %s", err)
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := ValidateConfig(); err != nil {
		logger.Errorf("Config validation failed, %v", err)
		return fmt.Errorf("config validation failed: %w", err)
	}

	debug := false
	if viper.IsSet("DEBUG") {
		debug = viper.GetBool("DEBUG")
	} else if viper.GetString("APP_ENV") == "" || viper.GetString("APP_ENV") == EnvDevelopment {
		debug = true
	}
	allowedHosts := viper.GetString("ALLOWED_HOSTS")
	if allowedHosts == "" {
		allowedHosts = "0.0.0.0"
	}
	rps := viper.GetInt("RATE_LIMIT_RPS")
	if rps <= 0 {
		rps = 100
	}
	burst := viper.GetInt("RATE_LIMIT_BURST")
	if burst <= 0 {
		burst = 200
	}
	timezone := viper.GetString("SERVER_TIMEZONE")
	if timezone == "" {
		timezone = "UTC"
	}
	requestTimeout := viper.GetInt("REQUEST_TIMEOUT_SECONDS")
	if requestTimeout <= 0 {
		requestTimeout = 30
	}
	shutdownTimeout := viper.GetInt("SERVER_SHUTDOWN_TIMEOUT")
	if shutdownTimeout <= 0 {
		shutdownTimeout = 10
	}
	refreshTokenTTLDays := viper.GetInt("REFRESH_TOKEN_TTL_DAYS")
	if refreshTokenTTLDays <= 0 {
		refreshTokenTTLDays = 7
	}
	accessTokenTTLMinutes := viper.GetInt("ACCESS_TOKEN_TTL_MINUTES")
	if accessTokenTTLMinutes <= 0 {
		accessTokenTTLMinutes = 15
	}
	masterSsl := viper.GetString("MASTER_SSL_MODE")
	if masterSsl == "" {
		masterSsl = "disable"
	}
	replicaSsl := viper.GetString("REPLICA_SSL_MODE")
	if replicaSsl == "" {
		replicaSsl = masterSsl
	}
	isDev := GetEnvironment() == EnvDevelopment
	corsOrigins := resolveCORSOrigins(parseCommaSeparated(viper.GetString("CORS_ALLOWED_ORIGINS")), isDev)
	trustedProxies := resolveTrustedProxies(viper.GetString("TRUSTED_PROXIES"), isDev)

	appConfig = &Configuration{
		Server: ServerConfiguration{
			Host:                   viper.GetString("SERVER_HOST"),
			Port:                   viper.GetString("SERVER_PORT"),
			JWTSecret:              viper.GetString("JWT_SECRET"),
			Debug:                  debug,
			AllowedHosts:           allowedHosts,
			Timezone:               timezone,
			RequestTimeoutSeconds:  requestTimeout,
			ShutdownTimeoutSeconds: shutdownTimeout,
			RefreshTokenTTLDays:    refreshTokenTTLDays,
			AccessTokenTTLMinutes:  accessTokenTTLMinutes,
			RateLimitRPS:           rps,
			RateLimitBurst:         burst,
			RateLimitUseUser:       viper.GetBool("RATE_LIMIT_USE_USER"),
			CORSAllowedOrigins:     corsOrigins,
			TrustedProxies:         trustedProxies,
		},
		Database: DatabaseConfiguration{
			Dbname:          viper.GetString("MASTER_DB_NAME"),
			Username:        viper.GetString("MASTER_DB_USER"),
			Password:        viper.GetString("MASTER_DB_PASSWORD"),
			Host:            viper.GetString("MASTER_DB_HOST"),
			Port:            viper.GetString("MASTER_DB_PORT"),
			LogMode:         viper.GetBool("MASTER_DB_LOG_MODE"),
			SslMode:         masterSsl,
			ReplicaDbname:   viper.GetString("REPLICA_DB_NAME"),
			ReplicaUsername: viper.GetString("REPLICA_DB_USER"),
			ReplicaPassword: viper.GetString("REPLICA_DB_PASSWORD"),
			ReplicaHost:     viper.GetString("REPLICA_DB_HOST"),
			ReplicaPort:     viper.GetString("REPLICA_DB_PORT"),
			ReplicaSslMode:  replicaSsl,
		},
	}
	if appConfig.Server.Port == "" {
		appConfig.Server.Port = "8000"
	}
	if appConfig.Server.Host == "" {
		appConfig.Server.Host = "0.0.0.0"
	}
	// Replica defaults to master when not set
	if appConfig.Database.ReplicaDbname == "" {
		appConfig.Database.ReplicaDbname = appConfig.Database.Dbname
		appConfig.Database.ReplicaUsername = appConfig.Database.Username
		appConfig.Database.ReplicaPassword = appConfig.Database.Password
		appConfig.Database.ReplicaHost = appConfig.Database.Host
		appConfig.Database.ReplicaPort = appConfig.Database.Port
		appConfig.Database.ReplicaSslMode = appConfig.Database.SslMode
	}

	logger.Infof("Configuration loaded and validated successfully")
	logger.Infof("CORS allowed origins: %v", appConfig.Server.CORSAllowedOrigins)
	return nil
}

// ValidateConfig validates that all required configuration values are present
// and not using insecure default values
func ValidateConfig() error {
	// Required configuration keys
	requiredKeys := []string{
		"JWT_SECRET",
		"SERVER_HOST",
		"SERVER_PORT",
		"MASTER_DB_NAME",
		"MASTER_DB_USER",
		"MASTER_DB_PASSWORD",
		"MASTER_DB_HOST",
		"MASTER_DB_PORT",
	}

	// Check all required keys are present
	var missingKeys []string
	for _, key := range requiredKeys {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing required config keys: %s", strings.Join(missingKeys, ", "))
	}

	// Validate JWT_SECRET is not using example value
	jwtSecret := viper.GetString("JWT_SECRET")
	if strings.Contains(strings.ToUpper(jwtSecret), "CHANGE-THIS") ||
		strings.Contains(strings.ToLower(jwtSecret), "your-jwt-secret") ||
		len(jwtSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be changed from example value and be at least 32 characters long. Generate with: openssl rand -base64 32")
	}

	// Validate port is numeric
	port := viper.GetString("SERVER_PORT")
	if port == "" {
		return fmt.Errorf("SERVER_PORT cannot be empty")
	}

	// Validate CORS origins (required in production; must be well-formed everywhere)
	rawOrigins := parseCommaSeparated(viper.GetString("CORS_ALLOWED_ORIGINS"))
	if err := validateCORSOrigins(rawOrigins, GetEnvironment() == EnvProduction); err != nil {
		return err
	}

	logger.Debugf("Config validation passed for all required keys")
	return nil
}

// GetString returns a string configuration value
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns an integer configuration value
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a boolean configuration value
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// Environment helper functions

// GetEnvironment returns the current application environment
// Defaults to "development" if APP_ENV is not set
func GetEnvironment() string {
	env := viper.GetString("APP_ENV")
	if env == "" {
		return EnvDevelopment
	}
	return env
}

// IsDevelopment returns true if running in development environment
func IsDevelopment() bool {
	return GetEnvironment() == EnvDevelopment
}

// IsStaging returns true if running in staging environment
func IsStaging() bool {
	return GetEnvironment() == EnvStaging
}

// IsProduction returns true if running in production environment
func IsProduction() bool {
	return GetEnvironment() == EnvProduction
}

// IsDebugEnabled returns true if debug mode is enabled
// Automatically returns true for development environment unless explicitly set to false
func IsDebugEnabled() bool {
	// If DEBUG is explicitly set, use that value
	if viper.IsSet("DEBUG") {
		return viper.GetBool("DEBUG")
	}

	// Otherwise, default to true for development, false for others
	return IsDevelopment()
}
