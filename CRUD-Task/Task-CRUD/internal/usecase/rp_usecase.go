package usecase

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/repository"
)

// RepoUseCase menangani logika bisnis untuk repository
type RepoUseCase struct {
	repoRepo repository.RepoRepositoryInterface
}

// NewRepoUseCase membuat instance baru dari RepoUseCase
func NewRepoUseCase(repoRepo repository.RepoRepositoryInterface) *RepoUseCase {
	return &RepoUseCase{repoRepo: repoRepo}
}

// GetAllRepos mengambil semua repository
func (uc *RepoUseCase) GetAllRepos() ([]entity.Repository, error) {
	return uc.repoRepo.GetAllRepos()
}

// GetRepositoryByID mengambil repository berdasarkan ID
func (uc *RepoUseCase) GetRepositoryByID(id uint) (*entity.Repository, error) {
	return uc.repoRepo.GetRepositoryByID(id)
}

// CreateRepo membuat repository baru
func (uc *RepoUseCase) CreateRepo(repo *entity.Repository) error {
	return uc.repoRepo.CreateRepo(repo)
}

// UpdateRepo memperbarui repository
func (uc *RepoUseCase) UpdateRepo(id uint, updatedRepo *entity.Repository) error {
	return uc.repoRepo.UpdateRepo(id, updatedRepo)
}

// DeleteRepo menghapus repository
func (uc *RepoUseCase) DeleteRepo(id uint) error {
	return uc.repoRepo.DeleteRepo(id)
}
