package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int    `db:"id"`
	FirstName      string `db:"first_name"`
	LastName       string `db:"last_name"`
	Email          string `db:"email"`
	HashedPassword []byte `db:"hashed_password"`
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
}
type UserModel struct {
	DB *sqlx.DB
}

func (m *UserModel) Insert(email, firstName, lastName, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (first_name, last_name, email, hashed_password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = m.DB.Exec(stmt, firstName, lastName, email, hashedPassword, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
