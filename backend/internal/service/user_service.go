package service

import (
	"context"
	"errors"
	"seagame/ticket/backend/internal/models/user"
	"seagame/ticket/backend/internal/repository"
	util "seagame/ticket/backend/utils"
)

type UserService interface {
	Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error)
	Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error)
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	GetAllUsers(ctx context.Context) ([]user.User, error)
	UpdateUser(ctx context.Context, user *user.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserServiceIml struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceIml{userRepo: userRepo}
}

func (s *UserServiceIml) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	existing, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err == nil && existing != nil {
		return nil, errors.New("user already exists")
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	role := req.Role
	if role == "" {
		role = "attendant"
	}

	newUser := &user.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     role,
	}
	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, errors.New("failed to create user")
	}
	return &user.RegisterResponse{
		ID:       newUser.ID.String(),
		Username: newUser.Username,
		Role:     newUser.Role,
	}, nil
}

func (s *UserServiceIml) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}
	foundUser, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := util.CheckPassword(req.Password, foundUser.Password); err != nil {
		return nil, errors.New("invalid password")
	}
	token, err := util.GenerateToken(foundUser.ID.String(), foundUser.Username, foundUser.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return &user.LoginResponse{
		Username: foundUser.Username,
		Token:    token,
		Role:     foundUser.Role,
	}, nil
}

func (s *UserServiceIml) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserServiceIml) GetAllUsers(ctx context.Context) ([]user.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserServiceIml) UpdateUser(ctx context.Context, user *user.User) error {
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *UserServiceIml) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.DeleteUser(ctx, id)
}
