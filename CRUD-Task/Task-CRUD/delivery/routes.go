package delivery

import (
	httpDelivery "Task-CRUD/delivery/http" // Alias untuk handler HTTP
	"Task-CRUD/internal/repository"
	"Task-CRUD/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewRouter mengatur semua rute aplikasi
func NewRouter(db *gorm.DB, redisClient *redis.Client) *mux.Router {
	router := mux.NewRouter()

	// Inisialisasi Repository dan Use Case
	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUseCase(userRepo)
	userHandler := httpDelivery.NewUserHandler(userUC)

	repoRepo := repository.NewRepoRepository(db)
	repoUC := usecase.NewRepoUseCase(repoRepo)
	repoHandler := httpDelivery.NewRepoHandler(repoUC)

	// User Routes
	router.Methods("GET").Path("/users").HandlerFunc(userHandler.GetUsers)
	router.Methods("POST").Path("/users").HandlerFunc(userHandler.CreateUser) // Tambahkan endpoint CreateUser

	// Repository Routes
	router.Methods("GET").Path("/repositories").HandlerFunc(repoHandler.GetRepositories)
	router.Methods("POST").Path("/repositories").HandlerFunc(repoHandler.CreateRepository)

	return router
}
