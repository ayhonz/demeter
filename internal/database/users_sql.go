package database

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserParams struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Storage) CreateUser(userParams CreateUserParams) (DBUser, error) {
	stmt := "INSERT INTO users (id, first_name, last_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING *"
	var user DBUser
	err := s.DB.Get(
		&user,
		stmt,
		userParams.ID,
		userParams.FirstName,
		userParams.LastName,
		userParams.CreatedAt,
		userParams.UpdatedAt,
	)
	if err != nil {
		return DBUser{}, err
	}

	return user, nil
}
