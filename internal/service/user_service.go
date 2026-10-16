package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lakejorge/boilerplate-golang/internal/domain"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := domain.ToUserResponse(user)
	return &res, nil
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]domain.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, total, err := s.userRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return domain.ToUserResponses(users), total, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != "" && req.Email != user.Email {
		existing, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && existing != nil && existing.ID != id {
			return nil, domain.ErrEmailAlreadyExists
		}
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	res := domain.ToUserResponse(user)
	return &res, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
