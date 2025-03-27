package repository

import (
	"Task-CRUD/internal/entity"
	"time"

	"gorm.io/gorm"
)

// UserRepositoryInterface mendefinisikan operasi CRUD untuk User
type UserRepositoryInterface interface {
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(id uint, user *entity.User) error
	DeleteUser(id uint) error
}

// UserRepository mengimplementasikan UserRepositoryInterface
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// GetAllUsers mengambil semua data user
func (r *UserRepository) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

// GetUserByID mengambil user berdasarkan ID
func (r *UserRepository) GetUserByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// CreateUser menambahkan pengguna baru ke database
func (r *UserRepository) CreateUser(user *entity.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return r.db.Create(user).Error
}

// UpdateUser memperbarui data user berdasarkan ID
func (r *UserRepository) UpdateUser(id uint, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Model(&entity.User{}).Where("id = ?", id).Updates(user).Error
}

// DeleteUser menghapus user berdasarkan ID
func (r *UserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}
