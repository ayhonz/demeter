package main

import (
	"fmt"
	"net/http"
	"strings"
)

type CookBookStore interface {
	GetRecipe(name string) string
	RecordRecipe(name string)
}

type CookBookServer struct {
	store CookBookStore
}

func (c *CookBookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recipe := strings.TrimPrefix(r.URL.Path, "/recipes/")
	switch r.Method {
	case http.MethodPost:
		c.processRecipe(w, recipe)
	case http.MethodGet:
		c.showRecipe(w, recipe)
	}
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
