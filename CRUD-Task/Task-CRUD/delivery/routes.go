package delivery

import (
	httpDelivery "Task-CRUD/delivery/http" // Alias untuk handler HTTP
	"Task-CRUD/internal/repository"
	"Task-CRUD/internal/usecase"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// NewRouter mengatur semua rute aplikasi (tanpa Redis)
func NewRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	// Inisialisasi Repository dan Use Case untuk User
	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUseCase(userRepo)         // ✅ Ini harus mengembalikan UserUseCaseInterface
	userHandler := httpDelivery.NewUserHandler(userUC) // ✅ Harus sesuai dengan parameter UserHandler

	// Inisialisasi Repository dan Use Case untuk Repository
	repoRepo := repository.NewRepoRepository(db) // ✅ Return interface
	repoUC := usecase.NewRepoUseCase(repoRepo)
	repoHandler := httpDelivery.NewRepoHandler(repoUC)

	// User Routes
	router.Methods("GET").Path("/users").HandlerFunc(userHandler.GetUsers)
	router.Methods("GET").Path("/users/{id}").HandlerFunc(userHandler.GetUserByID)
	router.Methods("POST").Path("/users").HandlerFunc(userHandler.CreateUser)
	router.Methods("PUT").Path("/users/{id}").HandlerFunc(userHandler.UpdateUser)
	router.Methods("DELETE").Path("/users/{id}").HandlerFunc(userHandler.DeleteUser)

	// Repository Routes (Update terbaru: Menggunakan GetAllRepos)
	router.Methods("GET").Path("/repositories").HandlerFunc(repoHandler.GetAllRepos) // ✅ FIXED
	router.Methods("GET").Path("/repositories/{id}").HandlerFunc(repoHandler.GetRepositoryByID)
	router.Methods("POST").Path("/repositories").HandlerFunc(repoHandler.CreateRepo)        // ✅ FIXED
	router.Methods("PUT").Path("/repositories/{id}").HandlerFunc(repoHandler.UpdateRepo)    // ✅ FIXED
	router.Methods("DELETE").Path("/repositories/{id}").HandlerFunc(repoHandler.DeleteRepo) // ✅ FIXED

	return router
}
