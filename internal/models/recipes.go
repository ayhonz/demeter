package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Recipe struct {
	ID          int            `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Ingredients pq.StringArray `db:"ingredients"`
	Categories  pq.StringArray `db:"categories"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	UserID      int            `db:"user_id"`
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

func (m *RecipeModel) Get(id string) (Recipe, error) {
	stmt := `SELECT * FROM recipes WHERE id=$1;`
	var recipe Recipe
	err := m.DB.Get(&recipe, stmt, id)
	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}

func (m *RecipeModel) List() ([]Recipe, error) {
	return []Recipe{}, nil
}
