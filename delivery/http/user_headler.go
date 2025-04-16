package http

import (
	"Task-CRUD/internal/entity"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	interfaces "Task-CRUD/internal/interfaces"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
)

type UserHandler struct {
	userUC interfaces.UserUseCaseInterface
}

func NewUserHandler(userUC interfaces.UserUseCaseInterface) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// --- Helper ---

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeSuccessResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func parseIDFromVars(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid ID")
	}
	return uint(id), nil
}

// --- Handlers ---

// GET /users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("Handler.GetUsers")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(r.Context(), span)

	users, err := h.userUC.GetUsers(ctx)
	if err != nil {
		log.Printf("ERROR | GetUsers: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Gagal mengambil data user")
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("Handler.GetUserByID")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(r.Context(), span)

	id, err := parseIDFromVars(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	log.Printf("INFO | GetUserByID ID: %d", id)

	user, err := h.userUC.GetUserByID(ctx, id)
	if err != nil {
		log.Printf("ERROR | GetUserByID: %v", err)
		writeErrorResponse(w, http.StatusNotFound, "User tidak ditemukan")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("Handler.CreateUser")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(r.Context(), span)

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	log.Printf("DEBUG | CreateUser payload: %+v", user)

	if err := h.userUC.CreateUser(ctx, &user); err != nil {
		log.Printf("ERROR | CreateUser: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusCreated, "User berhasil dibuat")
}

// PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("Handler.UpdateUser")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(r.Context(), span)

	id, err := parseIDFromVars(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	log.Printf("DEBUG | UpdateUser ID: %d Payload: %+v", id, user)

	if err := h.userUC.UpdateUser(ctx, id, &user); err != nil {
		log.Printf("ERROR | UpdateUser: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "User berhasil diperbarui")
}

// DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("Handler.DeleteUser")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(r.Context(), span)

	id, err := parseIDFromVars(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	log.Printf("INFO | DeleteUser ID: %d", id)

	if err := h.userUC.DeleteUser(ctx, id); err != nil {
		log.Printf("ERROR | DeleteUser: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Gagal menghapus user")
		return
	}

	w.WriteHeader(http.StatusNoContent) // ✅ No Content: tidak ada body
}
