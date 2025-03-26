package tests

// import (
// 	"Task-CRUD/delivery/http"
// 	"Task-CRUD/internal/usecase"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestGetRepositoryByID(t *testing.T) {
// 	mockRepo := &MockRepoRepository{}
// 	useCase := usecase.NewRepositoryUseCase(mockRepo)
// 	handler := http.NewRepositoryHandler(useCase)

// 	req := httptest.NewRequest("GET", "/repositories/1", nil)
// 	w := httptest.NewRecorder()

// 	handler.GetRepositoryByID(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// }
