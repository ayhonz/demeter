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
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), recipe)
}

func TestRetrievingUsers(t *testing.T) {
	store := NewInMemoryRecipeStore()
	server := NewCookBookServer(store)

	t.Run("get Users", func(t *testing.T) {
		response := httptest.NewRecorder()
		want := []User{
			{"Dennosuke"},
		}
		server.ServeHTTP(response, newGetUserRequest())
		got := getUsersFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertUsers(t, got, want)
	})

}
