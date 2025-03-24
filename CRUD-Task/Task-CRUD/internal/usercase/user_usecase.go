package usecase

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/repository"
)

type UserUseCase struct {
	UserRepo *repository.UserRepository
}

func NewUserUseCase(repo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{UserRepo: repo}
}

func (u *UserUseCase) RegisterUser(user entity.User) error {
	return u.UserRepo.CreateUser(user)
}

// Kurang Interface (fokus ke regis aja), + main.go
