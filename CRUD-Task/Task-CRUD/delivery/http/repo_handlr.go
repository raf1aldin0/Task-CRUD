package http

import (
	"encoding/json"
	"net/http"
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/usecase"
)

type RepoHandler struct {
	repoUC *usecase.RepoUseCase
}

func NewRepoHandler(repoUC *usecase.RepoUseCase) *RepoHandler {
	return &RepoHandler{repoUC: repoUC}
}

// Get all repositories
func (h *RepoHandler) GetRepositories(w http.ResponseWriter, r *http.Request) {
	repos, err := h.repoUC.GetRepositories()
	if err != nil {
		http.Error(w, "Failed to fetch repositories", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(repos)
}

// Create a new repository
func (h *RepoHandler) CreateRepository(w http.ResponseWriter, r *http.Request) {
	var repo entity.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.repoUC.CreateRepository(&repo); err != nil {
		http.Error(w, "Failed to create repository", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(repo)
}
