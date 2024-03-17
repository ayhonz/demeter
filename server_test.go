package cookbook

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ayhonz/racook/internal/database"
	"github.com/google/uuid"
)

func TestHealthCheck(t *testing.T) {
	store := StubRecipeStore{}
	server := NewCookBookServer(&store)

	request, _ := http.NewRequest(http.MethodGet, "/v1/healthz", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertStatus(t, response.Code, http.StatusOK)
	assertContentType(t, response, jsonContentType)
	assertResponseBody(t, response.Body.String(), "\"OK\"")
}

func TestGETRecipes(t *testing.T) {
	recipeID := uuid.New()
	UserID := uuid.New()
	store := StubRecipeStore{
		recipes: []database.Recipe{
			{
				ID:          recipeID,
				Title:       "chicken",
				Description: "delicios chicken",
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				UserID:      UserID,
			},
			{
				ID:          recipeID,
				Title:       "pasta",
				Description: "delicios pasta",
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				UserID:      UserID,
			},
		},
	}

	server := NewCookBookServer(&store)

	t.Run("returns chicken recipe", func(t *testing.T) {
		request := newGetRecipeRequest(recipeID)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getRecipeFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertRecipe(t, database.Recipe(got), store.recipes[0])
	})

	t.Run("returns 404 on not found recipe", func(t *testing.T) {
		request := newGetRecipeRequest(uuid.New())
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertStatus(t, got, want)
	})

	t.Run("returns list of recipes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/recipes", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := getRecipesFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertRecipes(t, got, store.recipes)
	})
}

func TestStoreRecipes(t *testing.T) {
	store := StubRecipeStore{
		recipes: []database.Recipe{},
	}

	server := NewCookBookServer(&store)

	t.Run("it records recipe when POST", func(t *testing.T) {

		json := []byte(`{"title": "chicken", "description": "delicious chicken"}`)
		request := newPostRecipeRequest(json)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.recipes) != 1 {
			t.Errorf("got %d calls to RecordRecipe want %d", len(store.recipes), 1)
		}

		if store.recipes[0].Title != "chicken" {
			t.Errorf("did not store correct recipe got %q want %q", store.recipes[0].Title, "chicken")
		}
	})
}
