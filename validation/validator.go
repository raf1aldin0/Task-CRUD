package validation

import (
	"fmt"
	"Task-CRUD/internal/entity"
	"github.com/go-playground/validator/v10"
)

// Reuse validator instance
var validate = validator.New()

// ValidateUser akan memvalidasi entitas User
func ValidateUser(user *entity.User) error {
	// Validasi dengan validator
	if err := validate.Struct(user); err != nil {
		// Mengembalikan error format yang lebih jelas
		return fmt.Errorf("user validation failed: %v", err)
	}
	return nil
}

// ValidateRepository akan memvalidasi entitas Repository
func ValidateRepository(repo *entity.Repository) error {
	// Validasi dengan validator
	if err := validate.Struct(repo); err != nil {
		// Mengembalikan error format yang lebih jelas
		return fmt.Errorf("repository validation failed: %v", err)
	}
	return nil
}
