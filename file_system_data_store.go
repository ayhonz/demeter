package main

import (
	"encoding/json"
	"os"
)

type FileSystemRecipeStore struct {
	database *json.Encoder
	recipes  []Recipe
}

func NewFileSystemRecipeStore(database *os.File) *FileSystemRecipeStore {
	database.Seek(0, 0)
	recipes, _ := NewRecipeList(database)
	return &FileSystemRecipeStore{
		database: json.NewEncoder(&tape{database}),
		recipes:  recipes,
	}

}

func (f *FileSystemRecipeStore) GetRecipeList() []Recipe {
	return f.recipes
}

func (f *FileSystemRecipeStore) GetRecipe(name string) *Recipe {
	var recipe Recipe

	for _, r := range f.recipes {
		if r.Title == name {
			recipe = r
		}
	}
	return &recipe
}

func (f *FileSystemRecipeStore) RecordRecipe(name string) {
	newRecipe := Recipe{Title: name, Description: "", Ingredients: []string{}}
	f.recipes = append(f.recipes, newRecipe)

	f.database.Encode(f.recipes)
}
