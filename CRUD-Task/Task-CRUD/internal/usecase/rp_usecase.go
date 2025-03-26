package usecase

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/repository"
)

type RepoUseCase struct {
	repoRepo *repository.RepoRepository
}

func NewRepoUseCase(repoRepo *repository.RepoRepository) *RepoUseCase {
	return &RepoUseCase{repoRepo: repoRepo}
}

func (uc *RepoUseCase) GetRepositories() ([]entity.Repository, error) {
	return uc.repoRepo.GetAllRepos()
}

func (uc *RepoUseCase) CreateRepository(repo *entity.Repository) error {
	return uc.repoRepo.CreateRepo(repo)
}
