package views

import "racook/internal/models"

type TemplateData struct {
	CurrentYear   int
	Authenticated bool
	Recipe        models.Recipe
	Recipes       []models.Recipe
}
