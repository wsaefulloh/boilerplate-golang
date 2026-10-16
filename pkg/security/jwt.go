package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lakejorge/boilerplate-golang/internal/domain"
)

type JWTService struct {
	secretKey            string
	expirationHours      int
	refreshExpirationDays int
}

func NewJWTService(secretKey string, expirationHours, refreshExpirationDays int) *JWTService {
	return &JWTService{
		secretKey:            secretKey,
		expirationHours:      expirationHours,
		refreshExpirationDays: refreshExpirationDays,
	}
}

func (j *JWTService) GenerateAccessToken(userID uuid.UUID, email, role string) (string, int64, error) {
	expirationTime := time.Now().Add(time.Duration(j.expirationHours) * time.Hour)
	expiresInSeconds := int64(time.Until(expirationTime).Seconds())

	claims := &domain.JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresInSeconds, nil
}

func (j *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.refreshExpirationDays) * 24 * time.Hour)

	claims := &domain.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateToken(tokenStr string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}
