package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Recipe struct {
	ID          int
	Title       string
	Description string
	Ingredients []string
	Categories  []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	userID      int
}

type RecipeModel struct {
	DB *sqlx.DB
}

func (m *RecipeModel) Insert(title, description string, ingredients, categories pq.StringArray) (int, error) {
	stmt := `INSERT INTO recipes (title, description, created_at, updated_at, categories, ingredients, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;`

	var id int64
	err := m.DB.Get(&id, stmt, title, description, time.Now(), time.Now(), categories, ingredients, 1)
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *RecipeModel) Get(id int) (Recipe, error) {
	return Recipe{}, nil
}

func (m *RecipeModel) List() ([]Recipe, error) {
	return []Recipe{}, nil
}
