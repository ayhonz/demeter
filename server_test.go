package cookbook

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

		got := getRecipeFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertRecipe(t, got, store.recipes["chicken"])
	})
	t.Run("returns pasta recipe", func(t *testing.T) {
		request := newGetRecipeRequest("pasta")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getRecipeFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertRecipe(t, got, store.recipes["pasta"])
	})

	t.Run("returns 404 on missing recipe", func(t *testing.T) {
		request := newGetRecipeRequest("not-there")
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
		want := []Recipe{
			store.recipes["chicken"],
			store.recipes["pasta"],
		}

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertRecipes(t, got, want)
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
