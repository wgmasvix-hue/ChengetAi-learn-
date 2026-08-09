package unit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	appauth "github.com/wgmasvix-hue/ChengetAi-learn-/internal/auth"
)

func TestPasswordHashing(t *testing.T) {
	hash, err := appauth.HashPassword("super-secret-password")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if !appauth.VerifyPassword(hash, "super-secret-password") {
		t.Fatal("expected password verification to succeed")
	}
	if appauth.VerifyPassword(hash, "wrong-password") {
		t.Fatal("expected password verification to fail")
	}
}

func TestJWTGenerationAndValidation(t *testing.T) {
	user := &appauth.User{ID: uuid.New(), Email: "user@example.com", Role: "learner"}
	secret := "abcdefghijklmnopqrstuvwxyz123456"
	token, err := appauth.GenerateJWT(user, secret, time.Hour)
	if err != nil {
		t.Fatalf("generate jwt: %v", err)
	}
	claims, err := appauth.ParseToken(token, secret)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.UserID != user.ID.String() || claims.Role != "learner" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestRoleValidation(t *testing.T) {
	if _, err := appauth.ValidatePublicRole("school_admin"); err == nil {
		t.Fatal("expected forbidden role to fail")
	}
	if role, err := appauth.ValidatePublicRole("teacher"); err != nil || role != "teacher" {
		t.Fatalf("expected teacher to pass, got role=%q err=%v", role, err)
	}
}
