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
	recipes     map[string]string
	recipeCalls []string
	users       []User
}

func (s *StubRecipeStore) GetRecipe(name string) string {
	return s.recipes[name]
}

func (s *StubRecipeStore) RecordRecipe(name string) {
	s.recipeCalls = append(s.recipeCalls, name)
}

func (s *StubRecipeStore) GetUsers() []User {
	return s.users
}

func TestGETRecipes(t *testing.T) {
	store := StubRecipeStore{
		recipes: map[string]string{
			"chicken": "chicken recipe",
			"pasta":   "pasta recipe",
		},
	}

	server := NewCookBookServer(&store)

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
	server := NewCookBookServer(&store)

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

func TestUsers(t *testing.T) {

	t.Run("it returns the user table as JSON", func(t *testing.T) {
		wantedUsers := []User{
			{"John"},
			{"Dennosuke"},
		}

		store := StubRecipeStore{nil, nil, wantedUsers}
		server := NewCookBookServer(&store)

		request := newGetUserRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getUsersFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertUsers(t, got, wantedUsers)
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
