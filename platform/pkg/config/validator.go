package config

import (
	"fmt"
	"os"
	"strings"
)

// ConfigValidator validates critical configuration at startup
type ConfigValidator struct {
	environment string
	errors      []string
	warnings    []string
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator(environment string) *ConfigValidator {
	return &ConfigValidator{
		environment: environment,
		errors:      []string{},
		warnings:    []string{},
	}
}

// ValidateProductionSecrets ensures production-grade secret requirements
func (cv *ConfigValidator) ValidateProductionSecrets() error {
	if cv.environment != "production" {
		return nil // Skip in non-production
	}

	// CRITICAL: These secrets MUST be set to strong values in production
	criticalSecrets := map[string]int{
		"JWT_SECRET":         32,
		"DB_PASSWORD":        12,
		"REDIS_PASSWORD":     12,
		"MINIO_ROOT_PASSWORD": 12,
		"TYPESENSE_API_KEY":  16,
	}

	for secret, minLength := range criticalSecrets {
		value := os.Getenv(secret)

		// Check if set
		if value == "" {
			cv.errors = append(cv.errors, fmt.Sprintf("CRITICAL: %s not set", secret))
			continue
		}

		// Check for default/weak values (common in development)
		if cv.isDefaultValue(secret, value) {
			cv.errors = append(cv.errors,
				fmt.Sprintf("CRITICAL: %s set to default/weak value (must be changed for production)", secret))
			continue
		}

		// Check minimum length
		if len(value) < minLength {
			cv.errors = append(cv.errors,
				fmt.Sprintf("CRITICAL: %s too short (minimum %d characters, got %d)",
					secret, minLength, len(value)))
		}
	}

	// CRITICAL: CORS cannot be wildcard in production
	cors := os.Getenv("CORS_ALLOWED_ORIGINS")
	if cors == "*" {
		cv.errors = append(cv.errors,
			"CRITICAL: CORS_ALLOWED_ORIGINS cannot be '*' in production - specify exact domains")
	}

	// CRITICAL: Database SSL must be required
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	if dbSSLMode == "" || dbSSLMode == "disable" {
		cv.errors = append(cv.errors,
			"CRITICAL: DB_SSL_MODE must be 'require' in production")
	}

	// CRITICAL: HTTPS must be enforced
	forceHTTPS := os.Getenv("FORCE_HTTPS")
	if forceHTTPS != "true" {
		cv.errors = append(cv.errors,
			"CRITICAL: FORCE_HTTPS must be 'true' in production")
	}

	return cv.getError()
}

// ValidateDevelopmentSecrets ensures development safety
func (cv *ConfigValidator) ValidateDevelopmentSecrets() {
	if cv.environment == "production" {
		return // These checks don't apply in production
	}

	// Just warn in development
	defaultValues := map[string]string{
		"JWT_SECRET":         "dev_jwt_secret_change_in_prod",
		"DB_PASSWORD":        "chengetai_dev_password_change_in_prod",
		"MINIO_ROOT_PASSWORD": "minioadmin",
	}

	for secret, defaultVal := range defaultValues {
		value := os.Getenv(secret)
		if value == defaultVal {
			cv.warnings = append(cv.warnings,
				fmt.Sprintf("WARNING: %s is set to default development value", secret))
		}
	}
}

// ValidateCORSConfiguration validates CORS settings
func (cv *ConfigValidator) ValidateCORSConfiguration() {
	cors := os.Getenv("CORS_ALLOWED_ORIGINS")

	if cors == "" {
		cv.warnings = append(cv.warnings, "WARNING: CORS_ALLOWED_ORIGINS not set, using defaults")
		return
	}

	// Parse origins
	origins := strings.Split(cors, ",")

	for _, origin := range origins {
		origin = strings.TrimSpace(origin)

		// Validate format
		if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
			cv.errors = append(cv.errors,
				fmt.Sprintf("Invalid CORS origin: %s (must start with http:// or https://)", origin))
		}

		// Warn about insecure origins in production
		if cv.environment == "production" && strings.HasPrefix(origin, "http://") {
			cv.warnings = append(cv.warnings,
				fmt.Sprintf("WARNING: Insecure CORS origin in production: %s (should use https://)", origin))
		}
	}
}

// ValidateDatabaseConfiguration validates database settings
func (cv *ConfigValidator) ValidateDatabaseConfiguration() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		cv.errors = append(cv.errors, "DB_HOST not set")
		return
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		cv.warnings = append(cv.warnings, "DB_PORT not set, using default 5432")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		cv.errors = append(cv.errors, "DB_NAME not set")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		cv.errors = append(cv.errors, "DB_USER not set")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		cv.errors = append(cv.errors, "DB_PASSWORD not set")
	}
}

// ValidateJWTConfiguration validates JWT settings
func (cv *ConfigValidator) ValidateJWTConfiguration() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		cv.errors = append(cv.errors, "JWT_SECRET not set")
		return
	}

	if len(jwtSecret) < 32 {
		cv.errors = append(cv.errors,
			fmt.Sprintf("JWT_SECRET too short (minimum 32 characters, got %d)", len(jwtSecret)))
	}

	jwtExpiry := os.Getenv("JWT_EXPIRY")
	if jwtExpiry == "" {
		cv.warnings = append(cv.warnings, "JWT_EXPIRY not set, using default 15m")
	}
}

// ValidateAllConfigurations runs all validation checks
func (cv *ConfigValidator) ValidateAllConfigurations() error {
	cv.ValidateProductionSecrets()
	cv.ValidateDevelopmentSecrets()
	cv.ValidateCORSConfiguration()
	cv.ValidateDatabaseConfiguration()
	cv.ValidateJWTConfiguration()

	// Print warnings
	for _, warning := range cv.warnings {
		fmt.Printf("⚠️  %s\n", warning)
	}

	// Return error if any critical issues
	return cv.getError()
}

// isDefaultValue checks if a secret is set to a known default
func (cv *ConfigValidator) isDefaultValue(secret, value string) bool {
	defaults := map[string][]string{
		"JWT_SECRET": {
			"dev_jwt_secret_change_in_prod",
			"dev-jwt-secret",
			"secret",
			"jwt-secret-key",
		},
		"DB_PASSWORD": {
			"chengetai_dev_password_change_in_prod",
			"password",
			"123456",
			"admin",
		},
		"MINIO_ROOT_PASSWORD": {
			"minioadmin",
			"admin",
		},
		"TYPESENSE_API_KEY": {
			"xyz",
			"test",
			"key",
		},
	}

	if defaultVals, exists := defaults[secret]; exists {
		for _, def := range defaultVals {
			if value == def {
				return true
			}
		}
	}

	return false
}

// getError returns error if any critical issues found
func (cv *ConfigValidator) getError() error {
	if len(cv.errors) > 0 {
		return fmt.Errorf("Configuration validation failed:\n- %s",
			strings.Join(cv.errors, "\n- "))
	}
	return nil
}

// HasErrors returns true if validation found critical errors
func (cv *ConfigValidator) HasErrors() bool {
	return len(cv.errors) > 0
}

// HasWarnings returns true if validation found warnings
func (cv *ConfigValidator) HasWarnings() bool {
	return len(cv.warnings) > 0
}

// PrintReport prints validation results
func (cv *ConfigValidator) PrintReport() {
	if len(cv.errors) > 0 {
		fmt.Println("\n❌ CRITICAL CONFIGURATION ERRORS:")
		for _, err := range cv.errors {
			fmt.Printf("   - %s\n", err)
		}
	}

	if len(cv.warnings) > 0 {
		fmt.Println("\n⚠️  CONFIGURATION WARNINGS:")
		for _, warning := range cv.warnings {
			fmt.Printf("   - %s\n", warning)
		}
	}

	if len(cv.errors) == 0 && len(cv.warnings) == 0 {
		fmt.Println("\n✅ Configuration validation passed")
	}
}
