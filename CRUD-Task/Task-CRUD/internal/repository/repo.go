package repository

import (
	"Task-CRUD/internal/entity"

	"gorm.io/gorm"
)

type RepoRepository struct {
	db *gorm.DB
}

func NewRepoRepository(db *gorm.DB) *RepoRepository {
	return &RepoRepository{db: db}
}

// Get all repositories
func (r *RepoRepository) GetAllRepos() ([]entity.Repository, error) {
	var repos []entity.Repository
	if err := r.db.Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

// Create a new repository
func (r *RepoRepository) CreateRepo(repo *entity.Repository) error {
	return r.db.Create(repo).Error
}
