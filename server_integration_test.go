package cookbook

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingRecipesAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "[]")
	defer cleanDatabase()
	store, err := NewFileSystemRecipeStore(database)
	if err != nil {
		t.Fatalf("didn't expect an error but got one %v", err)
	}
	server := NewCookBookServer(store)

	recipe := "chicken"

	server.ServeHTTP(httptest.NewRecorder(), newPostRecipeRequest(recipe))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetRecipeRequest(recipe))

	want := Recipe{
		Title:       recipe,
		Description: "",
		Ingredients: []string{},
	}
	got := getRecipeFromResponse(t, response.Body)

	assertStatus(t, response.Code, http.StatusOK)
	assertContentType(t, response, jsonContentType)
	assertRecipe(t, got, want)
}
