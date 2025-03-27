package repository

import (
	"Task-CRUD/internal/entity"

	"gorm.io/gorm"
)

// RepoRepositoryInterface mendefinisikan kontrak CRUD untuk repository
type RepoRepositoryInterface interface {
	GetAllRepos() ([]entity.Repository, error)
	GetRepositoryByID(id uint) (*entity.Repository, error)
	CreateRepo(repo *entity.Repository) error
	UpdateRepo(id uint, updatedRepo *entity.Repository) error
	DeleteRepo(id uint) error
}

// RepoRepository adalah implementasi dari RepoRepositoryInterface
type RepoRepository struct {
	db *gorm.DB
}

// NewRepoRepository membuat instance baru dari RepoRepository
func NewRepoRepository(db *gorm.DB) RepoRepositoryInterface {
	return &RepoRepository{db: db}
}

// GetAllRepos mengambil semua repository
func (r *RepoRepository) GetAllRepos() ([]entity.Repository, error) {
	var repos []entity.Repository
	if err := r.db.Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

// GetRepositoryByID mengambil repository berdasarkan ID
func (r *RepoRepository) GetRepositoryByID(id uint) (*entity.Repository, error) {
	var repo entity.Repository
	if err := r.db.First(&repo, id).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

// CreateRepo menambahkan repository baru
func (r *RepoRepository) CreateRepo(repo *entity.Repository) error {
	return r.db.Create(repo).Error
}

// UpdateRepo memperbarui data repository
func (r *RepoRepository) UpdateRepo(id uint, updatedRepo *entity.Repository) error {
	return r.db.Model(&entity.Repository{}).Where("id = ?", id).Updates(updatedRepo).Error
}

// DeleteRepo menghapus repository berdasarkan ID
func (r *RepoRepository) DeleteRepo(id uint) error {
	return r.db.Delete(&entity.Repository{}, id).Error
}
