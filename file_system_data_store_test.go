package cookbook

import (
	"os"
	"reflect"
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

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func assertRecipeEquals(t *testing.T, got, want Recipe) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("didn't expect an error but got one, %v", err)
	}
}
