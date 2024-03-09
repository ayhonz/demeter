package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingRecipesAndRetrivingThem(t *testing.T) {
	store := NewInMemoryRecipeStore()
	server := CookBookServer{store}
	recipe := "chicken"

	server.ServeHTTP(httptest.NewRecorder(), newPostRecipeRequest(recipe))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetRecipeRequest(recipe))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), recipe)
}
