package cookbook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/ayhonz/racook/internal/database"
	"github.com/google/uuid"
)

type StubRecipeStore struct {
	recipes []database.Recipe
	users   []database.User
}

func (s *StubRecipeStore) CreateRecipe(ctx context.Context, arg database.CreateRecipeParams) (database.Recipe, error) {
	recipe := database.Recipe{
		ID:          arg.ID,
		CreatedAt:   arg.CreatedAt,
		UpdatedAt:   arg.UpdatedAt,
		Title:       arg.Title,
		Description: arg.Description,
		UserID:      arg.UserID,
	}
	s.recipes = append(s.recipes, recipe)
	return recipe, nil
}

func (s *StubRecipeStore) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return database.User{
		ID:        arg.ID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
	}, nil
}

func (s *StubRecipeStore) GetRecipeByID(ctx context.Context, id uuid.UUID) (database.Recipe, error) {
	for _, recipe := range s.recipes {
		if recipe.ID == id {
			return recipe, nil
		}
	}
	return database.Recipe{}, fmt.Errorf("recipe not found")
}

func (s *StubRecipeStore) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}

	return database.User{}, fmt.Errorf("user not found")
}

func (s *StubRecipeStore) GetRecipes(ctx context.Context) ([]database.Recipe, error) {
	return s.recipes, nil
}

func getRecipesFromResponse(t testing.TB, body io.Reader) []Recipe {
	t.Helper()
	recipeList, err := NewRecipeList(body)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Recipes, '%v'", body, err)
	}

	return recipeList
}

func newGetRecipeRequest(id uuid.UUID) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/recipes/%s", id), nil)
	return req
}

func newPostRecipeRequest(json []byte) *http.Request {
	body := bytes.NewReader(json)
	req, _ := http.NewRequest(http.MethodPost, "/v1/recipes", body)
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

func assertRecipe(t testing.TB, got, want database.Recipe) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\n want %+v", got, want)
	}
}

func assertRecipes(t testing.TB, got []Recipe, want []database.Recipe) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("got %d recipes want %d", len(got), len(want))
	}
	for i, v := range got {
		if !reflect.DeepEqual(database.Recipe(v), want[i]) {
			t.Errorf("got %v want %v", got, want)
		}
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
