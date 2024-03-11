package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubRecipeStore struct {
	recipes     map[string]Recipe
	recipeCalls []string
	users       []User
}

func (s *StubRecipeStore) GetRecipe(name string) *Recipe {
	recipe, ok := s.recipes[name]
	if !ok {
		return nil
	}

	return &recipe
}

func (s *StubRecipeStore) RecordRecipe(name string) {
	s.recipeCalls = append(s.recipeCalls, name)
}

func (s *StubRecipeStore) GetUsers() []User {
	return s.users
}
func (s *StubRecipeStore) RecordUser(name string) {
	s.users = append(s.users, User{name})
}

func TestGETRecipes(t *testing.T) {
	store := StubRecipeStore{
		recipes: map[string]Recipe{
			"chicken": {
				Title:       "chicken recipe",
				Description: "",
				Ingredients: []string{},
			},
			"pasta": {
				Title:       "pasta recipe",
				Description: "",
				Ingredients: []string{},
			},
		},
	}

	server := NewCookBookServer(&store)

	t.Run("returns chicken recipe", func(t *testing.T) {
		request := newGetRecipeRequest("chicken")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		recipe := getRecipeFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertRecipe(t, recipe, store.recipes["chicken"])
	})
	t.Run("returns pasta recipe", func(t *testing.T) {
		request := newGetRecipeRequest("pasta")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		recipe := getRecipeFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertRecipe(t, recipe, store.recipes["pasta"])
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
		recipes: map[string]Recipe{},
	}
	server := NewCookBookServer(&store)

	t.Run("it records recipe when POST", func(t *testing.T) {
		recipe := "chicken"

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

func newGetUserRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	return request
}

func newPostUserRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s", name), nil)
	return request
}

func getRecipeFromResponse(t testing.TB, body io.Reader) (recipe Recipe) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&recipe)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of User, '%v'", body, err)
	}

	return
}

func getUsersFromResponse(t testing.TB, body io.Reader) (users []User) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&users)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of User, '%v'", body, err)
	}

	return
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

func assertRecipe(t testing.TB, got, want Recipe) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertUsers(t testing.TB, got, want []User) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
