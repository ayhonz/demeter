package cookbook

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileSystemRecipeStore struct {
	database *json.Encoder
	recipes  []Recipe
}

func NewFileSystemRecipeStore(file *os.File) (*FileSystemRecipeStore, error) {
	err := initialiseRecipeDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)

	}

	recipes, err := NewRecipeList(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading recipe store from file %s %v", file.Name(), err)

	}
	return &FileSystemRecipeStore{
		database: json.NewEncoder(&tape{file}),
		recipes:  recipes,
	}, nil

}

func initialiseRecipeDBFile(file *os.File) error {
	file.Seek(0, 0)
	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
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
