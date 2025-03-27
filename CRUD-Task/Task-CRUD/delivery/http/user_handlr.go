package http

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// UserHandler menangani request HTTP untuk user.
type UserHandler struct {
	userUC usecase.UserUseCaseInterface
}

// NewUserHandler membuat instance UserHandler.
func NewUserHandler(userUC usecase.UserUseCaseInterface) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// GetUsers mengambil semua user.
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUC.GetUsers()
	if err != nil {
		log.Println("Error fetching users:", err)
		http.Error(w, `{"error": "Gagal mengambil data user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID mengambil user berdasarkan ID.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, `{"error": "ID tidak valid"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userUC.GetUserByID(uint(id))
	if err != nil {
		log.Println("User not found:", err)
		http.Error(w, `{"error": "User tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser membuat user baru.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Format JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	// Validasi data
	if user.Name == "" || user.Email == "" {
		http.Error(w, `{"error": "Nama dan Email harus diisi"}`, http.StatusBadRequest)
		return
	}

	err := h.userUC.CreateUser(&user)
	if err != nil {
		log.Println("Failed to create user:", err)
		http.Error(w, `{"error": "Gagal membuat user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil dibuat"})
}

// UpdateUser memperbarui user berdasarkan ID.
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, `{"error": "ID tidak valid"}`, http.StatusBadRequest)
		return
	}

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Format JSON tidak valid"}`, http.StatusBadRequest)
		return
	}

	// Validasi data
	if user.Name == "" {
		http.Error(w, `{"error": "Nama tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	err = h.userUC.UpdateUser(uint(id), &user)
	if err != nil {
		log.Println("Failed to update user:", err)
		http.Error(w, `{"error": "Gagal memperbarui user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil diperbarui"})
}

// DeleteUser menghapus user berdasarkan ID.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, `{"error": "ID tidak valid"}`, http.StatusBadRequest)
		return
	}

	err = h.userUC.DeleteUser(uint(id))
	if err != nil {
		log.Println("Failed to delete user:", err)
		http.Error(w, `{"error": "Gagal menghapus user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, `{"message": "User berhasil dihapus"}`)
}
