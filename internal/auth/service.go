package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrForbiddenRole      = errors.New("role cannot self-register")
	ErrUserNotFound       = errors.New("user not found")
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type RegisterInput struct {
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Password          string `json:"password"`
	Role              string `json:"role"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	EducationLevel    string `json:"education_level"`
	Grade             string `json:"grade"`
	PreferredLanguage string `json:"preferred_language"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  PublicUser `json:"user"`
}

type PublicUser struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Role      string    `json:"role"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Service struct {
	repo     *Repository
	pool     *pgxpool.Pool
	secret   string
	tokenTTL time.Duration
}

func NewService(repo *Repository, pool *pgxpool.Pool, secret string) *Service {
	return &Service{repo: repo, pool: pool, secret: secret, tokenTTL: 24 * time.Hour}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	email, err := validateEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(input.Password)) < 8 {
		return nil, &ValidationError{Message: "Password must be at least 8 characters long"}
	}
	if strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.LastName) == "" {
		return nil, &ValidationError{Message: "First name and last name are required"}
	}

	role, err := ValidatePublicRole(input.Role)
	if err != nil {
		return nil, err
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &User{
		ID:           uuid.New(),
		Email:        email,
		Phone:        strings.TrimSpace(input.Phone),
		PasswordHash: hash,
		Role:         role,
		FirstName:    strings.TrimSpace(input.FirstName),
		LastName:     strings.TrimSpace(input.LastName),
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err := s.repo.CreateUserTx(ctx, tx, user); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	switch role {
	case "learner":
		if err := s.repo.CreateLearnerProfileTx(ctx, tx, user.ID, input.EducationLevel, input.Grade, input.PreferredLanguage); err != nil {
			return nil, fmt.Errorf("create learner profile: %w", err)
		}
	case "teacher":
		if err := s.repo.CreateTeacherProfileTx(ctx, tx, user.ID); err != nil {
			return nil, fmt.Errorf("create teacher profile: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	token, err := GenerateJWT(user, s.secret, s.tokenTTL)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	if err := s.repo.StoreSession(ctx, token, s.tokenTTL); err != nil {
		return nil, fmt.Errorf("store session: %w", err)
	}

	return &AuthResponse{Token: token, User: toPublicUser(user)}, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	email, err := validateEmail(input.Email)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if !VerifyPassword(user.PasswordHash, input.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := GenerateJWT(user, s.secret, s.tokenTTL)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	if err := s.repo.StoreSession(ctx, token, s.tokenTTL); err != nil {
		return nil, fmt.Errorf("store session: %w", err)
	}

	return &AuthResponse{Token: token, User: toPublicUser(user)}, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return &ValidationError{Message: "Authentication token is required"}
	}
	return s.repo.DeleteSession(ctx, token)
}

func (s *Service) Me(ctx context.Context, userID string) (*PublicUser, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return nil, &ValidationError{Message: "User identifier is invalid"}
	}
	user, err := s.repo.FindByID(ctx, parsed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	public := toPublicUser(user)
	return &public, nil
}

func ValidatePublicRole(role string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(role))
	switch normalized {
	case "learner", "teacher", "parent":
		return normalized, nil
	case "school_admin", "platform_admin":
		return "", ErrForbiddenRole
	default:
		return "", &ValidationError{Message: "Role must be learner, teacher, or parent"}
	}
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("v1$%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash)), nil
}

func VerifyPassword(encodedHash, password string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 3 || parts[0] != "v1" {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}
	actual := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, uint32(len(expected)))
	return subtle.ConstantTimeCompare(expected, actual) == 1
}

func GenerateJWT(user *User, secret string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := TokenClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidCredentials
	}
	return claims, nil
}

func validateEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return "", &ValidationError{Message: "Email address is required"}
	}
	if _, err := mail.ParseAddress(normalized); err != nil {
		return "", &ValidationError{Message: "Email address is invalid"}
	}
	return normalized, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func toPublicUser(user *User) PublicUser {
	return PublicUser{
		ID:        user.ID.String(),
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}
