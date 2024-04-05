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

type authenticateUser struct {
	ID             int    `db:"id"`
	HashedPassword []byte `db:"hashed_password"`
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
	var dbUser authenticateUser

	stmt := `SELECT id, hashed_password FROM users WHERE email=$1`

	err := m.DB.Get(&dbUser, stmt, email)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(dbUser.HashedPassword, []byte(password))
	if err != nil {
		return 0, err
	}

	return dbUser.ID, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
