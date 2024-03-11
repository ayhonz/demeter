package main

import (
	"encoding/json"
	"net/http"
)

const jsonContentType = "application/json"

type CookBookStore interface {
	GetRecipe(name string) *Recipe
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

	router := http.NewServeMux()

	router.Handle("GET /recipes", http.HandlerFunc(s.getRecipeHandler))
	router.Handle("GET /recipes/{recipeName}", http.HandlerFunc(s.getRecipeHandler))
	router.Handle("POST /recipes/{recipeName}", http.HandlerFunc(s.postRecipeHandler))

	s.Handler = router

	return s
}

func (c *CookBookServer) getRecipeHandler(w http.ResponseWriter, r *http.Request) {
	recipe := r.PathValue("recipeName")
	w.Header().Set("content-type", jsonContentType)

	c.showRecipe(w, recipe)
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
