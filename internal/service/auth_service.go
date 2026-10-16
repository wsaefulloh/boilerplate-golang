package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lakejorge/boilerplate-golang/internal/domain"
	"github.com/lakejorge/boilerplate-golang/pkg/security"
)

type authService struct {
	userRepo   domain.UserRepository
	jwtService *security.JWTService
}

func NewAuthService(userRepo domain.UserRepository, jwtService *security.JWTService) domain.AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	res := domain.ToUserResponse(newUser)
	return &res, nil
}

func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !security.CheckPasswordHash(req.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := domain.ToUserResponse(user)
	return &res, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.TokenResponse, error) {
	claims, err := s.jwtService.ValidateToken(refreshTokenStr)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}
