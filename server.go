package cookbook

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const jsonContentType = "application/json"

type CookBookStore interface {
	GetRecipe(name string) *Recipe
	GetRecipeList() []Recipe
	RecordRecipe(name string)
}

type CookBookServer struct {
	store CookBookStore
	http.Handler
}

type User struct {
	Name string
}

type Recipe struct {
	Title       string
	Description string
	Ingredients []string
}

func NewCookBookServer(store CookBookStore) *CookBookServer {
	s := new(CookBookServer)

	s.store = store

	router := chi.NewRouter()

	v1Router := chi.NewRouter()
	router.Mount("/v1", v1Router)

	v1Router.Get("/healthz", s.healthHandler)
	v1Router.Get("/recipes", s.getRecipesHandler)
	v1Router.Get("/recipes/{recipeName}", s.getRecipeHandler)
	v1Router.Post("/recipes/{recipeName}", s.postRecipeHandler)

	s.Handler = router

	return s
}

func (c *CookBookServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, http.StatusOK, "OK")
}

func (c *CookBookServer) getRecipeHandler(w http.ResponseWriter, r *http.Request) {
	recipe := r.PathValue("recipeName")
	w.Header().Set("content-type", jsonContentType)

	c.showRecipe(w, recipe)
}

func (c *CookBookServer) getRecipesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	recipes := c.store.GetRecipeList()
	if recipes == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
	} else {
		json.NewEncoder(w).Encode(recipes)
	}
}

func (c *CookBookServer) postRecipeHandler(w http.ResponseWriter, r *http.Request) {
	recipe := r.PathValue("recipeName")

	c.processRecipe(w, recipe)
}

func (c *CookBookServer) showRecipe(w http.ResponseWriter, recipeName string) {
	recipe := c.store.GetRecipe(recipeName)
	if recipe == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
	} else {
		json.NewEncoder(w).Encode(recipe)
	}
}

func (c *CookBookServer) processRecipe(w http.ResponseWriter, recipe string) {
	c.store.RecordRecipe(recipe)
	w.WriteHeader(http.StatusAccepted)
}
