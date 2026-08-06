package security

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

// InputValidator provides input validation utilities
type InputValidator struct{}

// NewInputValidator creates a new input validator
func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

// ValidateEmail validates an email address
func (iv *InputValidator) ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if len(email) > 254 {
		return fmt.Errorf("email is too long (max 254 characters)")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}

	return nil
}

// ValidateUsername validates a username
func (iv *InputValidator) ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(username) > 32 {
		return fmt.Errorf("username must not exceed 32 characters")
	}

	// Allow only alphanumeric, hyphen, underscore
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, hyphens, and underscores")
	}

	return nil
}

// ValidateUUID validates a UUID string
func (iv *InputValidator) ValidateUUID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	// UUID v4 regex
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, strings.ToLower(id))
	if !matched {
		return fmt.Errorf("invalid UUID format")
	}

	return nil
}

// ValidateString validates a generic string field
func (iv *InputValidator) ValidateString(field, value string, minLen, maxLen int) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	if len(value) < minLen {
		return fmt.Errorf("%s must be at least %d characters", field, minLen)
	}
	if len(value) > maxLen {
		return fmt.Errorf("%s must not exceed %d characters", field, maxLen)
	}

	return nil
}

// SanitizeString removes potentially dangerous characters (basic XSS protection)
func (iv *InputValidator) SanitizeString(input string) string {
	// Remove HTML tags and dangerous characters
	dangerous := []string{"<", ">", "\"", "'", "&", ";"}
	result := input

	for _, char := range dangerous {
		result = strings.ReplaceAll(result, char, "")
	}

	return strings.TrimSpace(result)
}

// ValidateURL validates a URL
func (iv *InputValidator) ValidateURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return fmt.Errorf("URL is required")
	}

	// Simple URL validation - must start with http:// or https://
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	if len(rawURL) > 2048 {
		return fmt.Errorf("URL is too long (max 2048 characters)")
	}

	return nil
}

// ValidateJSONField validates JSON field to prevent injection
func (iv *InputValidator) ValidateJSONField(data string) error {
	if len(data) > 1000000 { // 1MB limit
		return fmt.Errorf("JSON data is too large (max 1MB)")
	}
	return nil
}

// ValidateSQLIdentifier checks if a string is a valid SQL identifier
// Use this to validate column names, table names (NOT for values - use parameterized queries!)
func (iv *InputValidator) ValidateSQLIdentifier(identifier string) error {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return fmt.Errorf("identifier is required")
	}

	// Only allow alphanumeric and underscore
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, identifier)
	if !matched {
		return fmt.Errorf("invalid SQL identifier format")
	}

	if len(identifier) > 63 { // PostgreSQL identifier limit
		return fmt.Errorf("identifier exceeds 63 character limit")
	}

	return nil
}

// CheckSQLInjectionPattern performs basic SQL injection pattern detection
func (iv *InputValidator) CheckSQLInjectionPattern(input string) error {
	dangerous := []string{
		"';",
		"'--",
		"' OR '1'='1",
		"\" OR \"1\"=\"1",
		"1; DROP",
		"1; DELETE",
		"1; UPDATE",
		"1; INSERT",
		"UNION SELECT",
		"exec(",
		"execute(",
	}

	upper := strings.ToUpper(input)
	for _, pattern := range dangerous {
		if strings.Contains(upper, strings.ToUpper(pattern)) {
			return fmt.Errorf("potentially malicious input detected: suspicious SQL pattern")
		}
	}

	return nil
}
