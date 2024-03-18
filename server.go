package cookbook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ayhonz/racook/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const jsonContentType = "application/json"

type CookBookServer struct {
	store database.Querier
	http.Handler
}

func NewCookBookServer(store database.Querier) *CookBookServer {
	server := new(CookBookServer)

	server.store = store

	router := chi.NewRouter()

	v1Router := chi.NewRouter()
	router.Mount("/v1", v1Router)

	v1Router.Get("/healthz", server.healthHandler)

	v1Router.Post("/recipes", server.createRecipeHandler)
	v1Router.Get("/recipes", server.getRecipes)
	v1Router.Get("/recipes/{id}", server.getRecipeByIDHandler)

	v1Router.Post("/users", server.createUserHandler)

	server.Handler = router

	return server
}

func (c *CookBookServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, http.StatusOK, "OK")
}

func (c *CookBookServer) createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	params := new(parameters)
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing JSON %v", err))
	}

	userId, _ := uuid.Parse("68f4db84-355e-431f-b7b1-cf4688946e14")

	recipe, err := c.store.CreateRecipe(r.Context(), database.CreateRecipeParams{
		ID:          uuid.New(),
		Title:       params.Title,
		Description: params.Description,
		UpdatedAt:   time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
		UserID:      userId,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Could't create user %v", err))
		return
	}

	responseWithJSON(w, 202, databaseRecipeToRecipe(recipe))
}

func (c *CookBookServer) getRecipeByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	recipeID, err := uuid.Parse(id)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing ID %v", err))
		return
	}

	recipe, err := c.store.GetRecipeByID(r.Context(), recipeID)
	if err != nil {
		responseWithError(w, 404, "Recipe not found")
		return
	}

	responseWithJSON(w, 200, databaseRecipeToRecipe(recipe))
}

func (c *CookBookServer) getRecipes(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	var limit int32 = 10
	var offset int32 = 0
	if limitStr != "" {
		var err error
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		limit = int32(limitInt)
		if err != nil {
			responseWithError(w, 400, fmt.Sprint("Invalid limit"))
			return
		}
	}
	if offsetStr != "" {
		var err error
		offsetInt, err := strconv.ParseInt(offsetStr, 10, 32)
		offset = int32(offsetInt)
		if err != nil {
			responseWithError(w, 400, fmt.Sprint("Invalid offset"))
			return
		}
	}

	recipes, err := c.store.GetRecipes(r.Context(), database.GetRecipesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Printf("error getting recipes %v", err)
		// probably we should return 422 or something
		responseWithError(w, 404, fmt.Sprint("Recipes not found"))
		return
	}

	response := make([]Recipe, len(recipes))

	for i, recipe := range recipes {
		response[i] = databaseRecipeToRecipe(recipe)
	}

	responseWithJSON(w, 200, response)
}

func (c *CookBookServer) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	params := new(parameters)
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing JSON %v", err))
	}
	user, err := c.store.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		FirstName: params.FirstName,
		LastName:  params.LastName,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Could't create user %v", err))
		return

	}

	responseWithJSON(w, 201, databaseUserToUser(user))
}
