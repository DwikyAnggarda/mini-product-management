package repository

import (
	"context"
	"database/sql"
	"fmt"

	"product-management/backend/internal/model"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = $1 LIMIT 1`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, fmt.Errorf("user not found")
		}
		return model.User{}, err
	}
	return user, nil
}
