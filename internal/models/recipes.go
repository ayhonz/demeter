package models

import (
	"time"

	"github.com/jmoiron/sqlx"
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

func (m *RecipeModel) Insert(title, description string, ingredients, categories []string) (int, error) {
	return 0, nil
}

func (m *RecipeModel) Get(id int) (Recipe, error) {
	return Recipe{}, nil
}

func (m *RecipeModel) List() ([]Recipe, error) {
	return []Recipe{}, nil
}
