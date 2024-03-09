package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const jsonContentType = "application/json"

type CookBookStore interface {
	GetRecipe(name string) string
	RecordRecipe(name string)
	GetUsers() []User
}

type CookBookServer struct {
	store CookBookStore
	http.Handler
}

type User struct {
	Name string
}

func NewCookBookServer(store CookBookStore) *CookBookServer {
	s := new(CookBookServer)

	s.store = store

	router := http.NewServeMux()

	router.Handle("GET /recipes/{recipeName}", http.HandlerFunc(s.getRecipeHandler))
	router.Handle("POST /recipes/{recipeName}", http.HandlerFunc(s.PostRecipeHandler))
	router.Handle("GET /users", http.HandlerFunc(s.getUserHandler))

	s.Handler = router

	return s
}

func (c *CookBookServer) getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(c.store.GetUsers())
}

func (c *CookBookServer) getUserTable() []User {
	return []User{
		{"John"},
		{"Dennosuke"},
	}

}

func (c *CookBookServer) getRecipeHandler(w http.ResponseWriter, r *http.Request) {
	recipe := r.PathValue("recipeName")

	c.showRecipe(w, recipe)
}

func (c *CookBookServer) PostRecipeHandler(w http.ResponseWriter, r *http.Request) {
	recipe := r.PathValue("recipeName")

	c.processRecipe(w, recipe)
}

func (c *CookBookServer) showRecipe(w http.ResponseWriter, recipe string) {
	recipeName := c.store.GetRecipe(recipe)
	if recipeName == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, recipeName)
}

func (c *CookBookServer) processRecipe(w http.ResponseWriter, recipe string) {
	c.store.RecordRecipe(recipe)
	w.WriteHeader(http.StatusAccepted)
}
