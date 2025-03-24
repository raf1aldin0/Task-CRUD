package repository

import (
	"database/sql"
	"CRUD TASK/Task-CRUD/internal/entity"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user entity.User) error {
	query := "INSERT INTO users (username, email) VALUES ($1, $2)"
	_, err := r.DB.Exec(query, user.Username, user.Email)
	return err
}

func (r *UserRepository) GetUserByID(id int) (*entity.User, error) {
	query := "SELECT id, username, email FROM users WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
