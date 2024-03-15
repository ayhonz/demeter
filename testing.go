package cookbook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

func (s *StubRecipeStore) GetRecipeList() []Recipe {
	var recipes []Recipe
	for _, recipe := range s.recipes {
		recipes = append(recipes, recipe)
	}

	return recipes
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

func getRecipesFromResponse(t testing.TB, body io.Reader) []Recipe {
	t.Helper()
	recipeList, err := NewRecipeList(body)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of User, '%v'", body, err)
	}

	return recipeList
}

func newGetRecipeRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/recipes/%s", name), nil)
	return req
}

func newPostRecipeRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/recipes/%s", name), nil)
	return req
}

func newGetUserRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/v1/users", nil)
	return request
}

func newPostUserRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/users/%s", name), nil)
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

func assertRecipes(t testing.TB, got, want []Recipe) {
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
