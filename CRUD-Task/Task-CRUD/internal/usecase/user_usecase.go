package usecase

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/repository"
	"time"
)

type UserUseCase struct {
	userRepo *repository.UserRepository
}

func NewUserUseCase(userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// GetUsers mengambil semua user dari database
func (uc *UserUseCase) GetUsers() ([]entity.User, error) {
	return uc.userRepo.GetAllUsers()
}

// CreateUser membuat user baru di database
func (uc *UserUseCase) CreateUser(user *entity.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return uc.userRepo.CreateUser(user)
}
