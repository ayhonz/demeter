package cookbook

import (
	"time"

	"github.com/ayhonz/racook/internal/database"
)

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type Recipe struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
	Categories  []string  `json:"categories"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
	}
}

func databaseRecipeToRecipe(dbRecipe database.Recipe) Recipe {
	return Recipe{
		ID:          dbRecipe.ID,
		CreatedAt:   dbRecipe.CreatedAt,
		UpdatedAt:   dbRecipe.UpdatedAt,
		Description: dbRecipe.Description,
		Title:       dbRecipe.Title,
		UserID:      dbRecipe.UserID,
		Categories:  dbRecipe.Categories,
	}
}
