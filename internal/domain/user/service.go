package user

import (
	"fmt"
	"time"

	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type Service interface {
	RegisterUser(req dto.RegisterRequest) (*dto.UserResponse, error)
	LoginUser(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) Service {
	return &service{repo: repo, jwtService: jwtService}
}

func (s *service) RegisterUser(req dto.RegisterRequest) (*dto.UserResponse, error) {
	role := req.Role
	if role == "" {
		role = RoleDriver
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  role,
	}

	// পাসওয়ার্ড hash করে user.Password-এ বসাও
	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.checkPassword(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
