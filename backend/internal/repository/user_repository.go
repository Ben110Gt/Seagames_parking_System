package repository

import (
	"seagame/ticket/backend/internal/models/user"

	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *user.User) error
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetAllUsers(ctx context.Context) ([]user.User, error)
	UpdateUser(ctx context.Context, user *user.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return &u, err
}

func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]user.User, error) {
	var users []user.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&user.User{}).Error
}

func (r *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var user user.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}
