package http

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// RepoHandler menangani HTTP request untuk repository
type RepoHandler struct {
	repoUC *usecase.RepoUseCase
}

// NewRepoHandler membuat instance baru dari RepoHandler
func NewRepoHandler(repoUC *usecase.RepoUseCase) *RepoHandler {
	return &RepoHandler{repoUC: repoUC}
}

// GetAllRepos menangani request untuk mendapatkan semua repository
func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := h.repoUC.GetAllRepos()
	if err != nil {
		http.Error(w, "Gagal mengambil daftar repository", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

// GetRepositoryByID menangani request untuk mendapatkan satu repository berdasarkan ID
func (h *RepoHandler) GetRepositoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	repo, err := h.repoUC.GetRepositoryByID(uint(id))
	if err != nil {
		http.Error(w, "Repository tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repo)
}

// CreateRepo menangani request untuk membuat repository baru
func (h *RepoHandler) CreateRepo(w http.ResponseWriter, r *http.Request) {
	var repo entity.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.repoUC.CreateRepo(&repo); err != nil {
		http.Error(w, "Gagal membuat repository", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(repo)
}

// UpdateRepo menangani request untuk memperbarui repository berdasarkan ID
func (h *RepoHandler) UpdateRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var updatedRepo entity.Repository
	if err := json.NewDecoder(r.Body).Decode(&updatedRepo); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.repoUC.UpdateRepo(uint(id), &updatedRepo); err != nil {
		http.Error(w, "Gagal memperbarui repository", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedRepo)
}

// DeleteRepo menangani request untuk menghapus repository berdasarkan ID
func (h *RepoHandler) DeleteRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.repoUC.DeleteRepo(uint(id)); err != nil {
		http.Error(w, "Gagal menghapus repository", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
