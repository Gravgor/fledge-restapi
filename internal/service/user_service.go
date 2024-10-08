package service

import (
	"context"
	"fledge-restapi/internal/domain/entity"
	"fledge-restapi/internal/domain/repository"
	"fledge-restapi/internal/util"
	"fledge-restapi/pkg/errors"
)

type UserService interface {
	CreateUser(ctx context.Context, req *entity.SignupRequest) error
	Login(ctx context.Context, req *entity.LoginRequest) (string, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *entity.SignupRequest) error {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return errors.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	return s.userRepo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, req *entity.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	if err := util.CheckPassword(req.Password, user.Password); err != nil {
		return "", errors.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := util.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return s.userRepo.FindByID(ctx, id)
}
