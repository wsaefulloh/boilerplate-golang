package security_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/wsaefulloh/boilerplate-golang/pkg/security"
)

func TestHashAndPassword(t *testing.T) {
	password := "mySecretPassword123"

	hash, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got: %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash string")
	}

	if !security.CheckPasswordHash(password, hash) {
		t.Fatal("expected password check to succeed")
	}

	if security.CheckPasswordHash("wrongPassword", hash) {
		t.Fatal("expected password check to fail for wrong password")
	}
}

func TestJWTService(t *testing.T) {
	secret := "testSecretKey123"
	jwtSvc := security.NewJWTService(secret, 1, 7)

	userID := uuid.New()
	email := "test@example.com"
	role := "user"

	// Test Access Token Generation & Validation
	accessToken, expiresIn, err := jwtSvc.GenerateAccessToken(userID, email, role)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	if accessToken == "" {
		t.Fatal("access token should not be empty")
	}

	if expiresIn <= 0 {
		t.Fatalf("expected positive expires_in, got %d", expiresIn)
	}

	claims, err := jwtSvc.ValidateToken(accessToken)
	if err != nil {
		t.Fatalf("failed to validate access token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}

	// Test Refresh Token
	refreshToken, err := jwtSvc.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	refClaims, err := jwtSvc.ValidateToken(refreshToken)
	if err != nil {
		t.Fatalf("failed to validate refresh token: %v", err)
	}

	if refClaims.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, refClaims.UserID)
	}
}
