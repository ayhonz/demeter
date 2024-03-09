package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubRecipeStore struct {
	recipes     map[string]string
	recipeCalls []string
}

func (s *StubRecipeStore) GetRecipe(name string) string {
	return s.recipes[name]
}

func (s *StubRecipeStore) RecordRecipe(name string) {
	s.recipeCalls = append(s.recipeCalls, name)
}

func TestGETRecipes(t *testing.T) {
	store := StubRecipeStore{
		recipes: map[string]string{
			"chicken": "chicken recipe",
			"pasta":   "pasta recipe",
		},
	}

	server := &CookBookServer{&store}

	t.Run("returns chicken recipe", func(t *testing.T) {
		request := newGetRecipeRequest("chicken")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "chicken recipe")
	})
	t.Run("returns pasta recipe", func(t *testing.T) {
		request := newGetRecipeRequest("pasta")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "pasta recipe")

	})
	t.Run("returns 404 on missing recipe", func(t *testing.T) {
		request := newGetRecipeRequest("not-there")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertStatus(t, got, want)

	})
}

func TestStoreRecipes(t *testing.T) {
	store := StubRecipeStore{
		recipes: map[string]string{},
	}
	server := &CookBookServer{&store}

	t.Run("it records recipe when POST", func(t *testing.T) {
		recipe := "turkey"

		request := newPostRecipeRequest(recipe)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.recipeCalls) != 1 {
			t.Errorf("got %d calls to RecordRecipe want %d", len(store.recipeCalls), 1)
		}
		if store.recipeCalls[0] != recipe {
			t.Errorf("did not store correct recipe got %q want %q", store.recipeCalls[0], recipe)
		}
	})
}

func newGetRecipeRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/recipes/%s", name), nil)
	return req
}

func newPostRecipeRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/recipes/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("status code is wrong, got %d, want %d", got, want)
	}
}
