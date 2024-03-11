package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingRecipesAndRetrievingThem(t *testing.T) {
	store := NewInMemoryRecipeStore()
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
