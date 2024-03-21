package database

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"hashed_password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateUserParams struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Storage) CreateUser(userParams CreateUserParams) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userParams.Password), 12)
	if err != nil {
		return User{}, err
	}

	stmt := "INSERT INTO users (first_name, last_name, email, hashed_password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *"
	var user User
	err = s.DB.Get(
		&user,
		stmt,
		userParams.FirstName,
		userParams.LastName,
		userParams.Email,
		hashedPassword,
		userParams.CreatedAt,
		userParams.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
