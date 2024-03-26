package database

import (
	"time"

	"github.com/lib/pq"
)

type Recipe struct {
	ID           int
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Categories   pq.StringArray `db:"categories"`
	Ingerediants pq.StringArray `db:"ingerediants"`
	UserID       int            `db:"user_id"`
}

type CreateRecipeParams struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description string
	Categories  pq.StringArray
	UserID      int
}

func (s *Storage) GetRecipeByID(id int) (Recipe, error) {
	var recipe Recipe
	err := s.DB.Get(&recipe, "SELECT * FROM recipes WHERE id = $1", id)

	return recipe, err
}

func (s *Storage) GetRecipes() ([]Recipe, error) {
	var recipes []Recipe
	err := s.DB.Select(&recipes, "SELECT * FROM recipes")
	return recipes, err
}

func (s *Storage) CreateRecipe(recipeParams CreateRecipeParams) (Recipe, error) {
	stmt := "INSERT INTO recipes (title, description, created_at, updated_at, categories, user_id) VALUES ($2, $2, $3, $4, $5, $6) RETURNING *"
	var recipe Recipe
	err := s.DB.Get(
		&recipe,
		stmt,
		recipeParams.Title,
		recipeParams.Description,
		recipeParams.CreatedAt,
		recipeParams.UpdatedAt,
		recipeParams.Categories,
		recipeParams.UserID,
	)
	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}
