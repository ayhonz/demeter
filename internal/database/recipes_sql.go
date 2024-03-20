package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CreateRecipeParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description string
	Categories  pq.StringArray
	UserID      uuid.UUID
}

func (s *Storage) GetRecipeByID(id uuid.UUID) (DBRecipe, error) {
	var recipe DBRecipe
	err := s.DB.Get(&recipe, "SELECT * FROM recipes WHERE id = $1", id)

	return recipe, err
}

func (s *Storage) GetRecipes() ([]DBRecipe, error) {
	var recipes []DBRecipe
	err := s.DB.Select(&recipes, "SELECT * FROM recipes")
	return recipes, err
}

func (s *Storage) CreateRecipe(recipeParams CreateRecipeParams) (DBRecipe, error) {
	stmt := "INSERT INTO recipes (id, title, description, created_at, updated_at, categories, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *"
	var recipe DBRecipe
	err := s.DB.Get(
		&recipe,
		stmt,
		recipeParams.ID,
		recipeParams.Title,
		recipeParams.Description,
		recipeParams.CreatedAt,
		recipeParams.UpdatedAt,
		recipeParams.Categories,
		recipeParams.UserID,
	)
	if err != nil {
		return DBRecipe{}, err
	}

	return recipe, nil
}
