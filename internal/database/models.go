package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DBUser struct {
	ID        uuid.UUID
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type DBRecipe struct {
	ID          uuid.UUID
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Categories  pq.StringArray `db:"categories"`
	UserID      uuid.UUID      `db:"user_id"`
}
