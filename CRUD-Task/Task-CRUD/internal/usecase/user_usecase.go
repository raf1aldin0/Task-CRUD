package usecase

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/repository"
)

// UserUseCaseInterface mendefinisikan operasi untuk User
type UserUseCaseInterface interface {
	GetUsers() ([]entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(id uint, user *entity.User) error
	DeleteUser(id uint) error
}

// UserUseCase mengimplementasikan business logic
type UserUseCase struct {
	userRepo repository.UserRepositoryInterface
}

// NewUserUseCase membuat instance baru UserUseCase
func NewUserUseCase(userRepo repository.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{userRepo: userRepo}
}

// GetUsers mengambil semua user
func (uc *UserUseCase) GetUsers() ([]entity.User, error) {
	return uc.userRepo.GetAllUsers()
}

// GetUserByID mengambil user berdasarkan ID
func (uc *UserUseCase) GetUserByID(id uint) (*entity.User, error) {
	return uc.userRepo.GetUserByID(id)
}

// CreateUser membuat user baru
func (uc *UserUseCase) CreateUser(user *entity.User) error {
	return uc.userRepo.CreateUser(user)
}

// UpdateUser memperbarui data user
func (uc *UserUseCase) UpdateUser(id uint, user *entity.User) error {
	return uc.userRepo.UpdateUser(id, user)
}

// DeleteUser menghapus user berdasarkan ID
func (uc *UserUseCase) DeleteUser(id uint) error {
	return uc.userRepo.DeleteUser(id)
}
