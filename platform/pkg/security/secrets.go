package security

import (
	"fmt"
	"os"
)

// SecretsManager handles secure retrieval of secrets from environment
type SecretsManager struct{}

// NewSecretsManager creates a new secrets manager
func NewSecretsManager() *SecretsManager {
	return &SecretsManager{}
}

// GetRequired retrieves a required environment variable
// Panics if not found
func (sm *SecretsManager) GetRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("FATAL: Required environment variable not set: %s", key))
	}
	if value == "" {
		panic(fmt.Sprintf("FATAL: Required environment variable is empty: %s", key))
	}
	return value
}

// GetOptional retrieves an optional environment variable with default
func (sm *SecretsManager) GetOptional(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

// ValidateSecrets performs validation of critical secrets
func (sm *SecretsManager) ValidateSecrets() error {
	requiredSecrets := []string{
		"DB_PASSWORD",
		"JWT_SECRET",
	}

	for _, secret := range requiredSecrets {
		value := sm.GetOptional(secret, "")
		if value == "" {
			return fmt.Errorf("critical secret not configured: %s", secret)
		}

		// Validate minimum length for sensitive values
		switch secret {
		case "JWT_SECRET":
			if len(value) < 32 {
				return fmt.Errorf("%s must be at least 32 characters long for security", secret)
			}
		case "DB_PASSWORD":
			if len(value) < 12 {
				return fmt.Errorf("%s must be at least 12 characters long for security", secret)
			}
		}
	}

	return nil
}

// GetDatabaseURL constructs a secure database connection URL
func (sm *SecretsManager) GetDatabaseURL() string {
	dbUser := sm.GetRequired("DB_USER")
	dbPassword := sm.GetRequired("DB_PASSWORD")
	dbHost := sm.GetRequired("DB_HOST")
	dbPort := sm.GetOptional("DB_PORT", "5432")
	dbName := sm.GetRequired("DB_NAME")

	// Note: Password is embedded here, consider using connection string parameter
	// or database connection pooling with secrets management
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require&statement_cache_mode=describe",
		dbUser, dbPassword, dbHost, dbPort, dbName)
}

// GetDatabaseURLSSLDisabled constructs a database URL for development (SSL disabled)
// ONLY for development/testing - should never be used in production
func (sm *SecretsManager) GetDatabaseURLSSLDisabled() string {
	dbUser := sm.GetRequired("DB_USER")
	dbPassword := sm.GetRequired("DB_PASSWORD")
	dbHost := sm.GetRequired("DB_HOST")
	dbPort := sm.GetOptional("DB_PORT", "5432")
	dbName := sm.GetRequired("DB_NAME")

	// WARNING: sslmode=disable should never be used in production
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
}

// GetJWTSecret retrieves and validates the JWT secret
func (sm *SecretsManager) GetJWTSecret() ([]byte, error) {
	secret := sm.GetRequired("JWT_SECRET")
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}
	return []byte(secret), nil
}

// GetAPIKey retrieves an API key for external services
func (sm *SecretsManager) GetAPIKey(service string) string {
	key := os.Getenv(service + "_API_KEY")
	if key == "" {
		// Return a placeholder that makes it obvious the key is not set
		return "not-configured"
	}
	return key
}

// VerifyProductionSecrets ensures production-grade security
func (sm *SecretsManager) VerifyProductionSecrets() error {
	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		return nil // Skip in non-production
	}

	// Production-only validations
	checks := []struct {
		env string
		min int
	}{
		{"JWT_SECRET", 32},
		{"DB_PASSWORD", 16},
		{"REDIS_PASSWORD", 16},
	}

	for _, check := range checks {
		value := os.Getenv(check.env)
		if value == "" {
			return fmt.Errorf("production: missing required secret: %s", check.env)
		}
		if len(value) < check.min {
			return fmt.Errorf("production: %s too short (minimum %d characters)", check.env, check.min)
		}
	}

	// Verify CORS is not set to wildcard in production
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins == "*" {
		return fmt.Errorf("production: CORS_ALLOWED_ORIGINS cannot be '*' - specify exact domains")
	}

	return nil
}
