package repository

import (
	"Task-CRUD/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAllUsers mengambil semua data pengguna
func (r *UserRepository) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser menambahkan pengguna baru ke database
func (r *UserRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}
