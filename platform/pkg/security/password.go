package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher provides secure password hashing using bcrypt
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher creates a new password hasher with bcrypt cost factor
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost < bcrypt.DefaultCost {
		cost = bcrypt.DefaultCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	return &PasswordHasher{cost: cost}
}

// Hash generates a secure bcrypt hash of the password
func (ph *PasswordHasher) Hash(password string) (string, error) {
	if len(password) < 8 {
		return "", fmt.Errorf("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		return "", fmt.Errorf("password must not exceed 72 characters (bcrypt limitation)")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// Verify checks if a password matches its hash
func (ph *PasswordHasher) Verify(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsValidPassword performs client-side password validation
func IsValidPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 72 {
		return fmt.Errorf("password must not exceed 72 characters")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		case r == '!' || r == '@' || r == '#' || r == '$' || r == '%' || r == '&' || r == '*':
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password must contain uppercase, lowercase, and numbers")
	}

	return nil
}
