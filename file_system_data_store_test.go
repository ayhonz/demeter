package cookbook

import (
	"testing"
)

func TestFileSystemRecipes(t *testing.T) {
	t.Run("Recipes from Reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
            {"Title": "Chicken recipe", "Description": "", "Ingredients": []},
            {"Title": "Pasta recipe", "Description": "", "Ingredients": []}
            ]`)
		defer cleanDatabase()

		store, err := NewFileSystemRecipeStore(database)
		assertNoError(t, err)

		got := store.GetRecipeList()
		want := []Recipe{
			{Title: "Chicken recipe", Description: "", Ingredients: []string{}},
			{Title: "Pasta recipe", Description: "", Ingredients: []string{}},
		}
		assertRecipes(t, got, want)

		got = store.GetRecipeList()

		assertRecipes(t, got, want)
	})

	t.Run("get recipe", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
            {"Title": "Chicken recipe", "Description": "", "Ingredients": []},
            {"Title": "Pasta recipe", "Description": "", "Ingredients": []}
            ]`)
		defer cleanDatabase()
		store, err := NewFileSystemRecipeStore(database)
		assertNoError(t, err)

		got := store.GetRecipe("Chicken recipe")
		want := Recipe{Title: "Chicken recipe", Description: "", Ingredients: []string{}}

		assertRecipeEquals(t, *got, want)
	})

	t.Run("store recipe", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
            ]`)
		defer cleanDatabase()

		store, err := NewFileSystemRecipeStore(database)
		assertNoError(t, err)
		store.RecordRecipe("Chicken recipe")

		got := store.GetRecipe("Chicken recipe")
		want := Recipe{Title: "Chicken recipe", Description: "", Ingredients: []string{}}

		assertRecipeEquals(t, *got, want)
	})

	t.Run("works with empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemRecipeStore(database)

		assertNoError(t, err)
	})
}
