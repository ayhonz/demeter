package main

import "sync"

type InMemoryRecipeStore struct {
	recipesStore map[string]Recipe
	usersStore   []User
	lock         sync.RWMutex
}

func NewInMemoryRecipeStore() *InMemoryRecipeStore {
	return &InMemoryRecipeStore{
		map[string]Recipe{},
		[]User{},
		sync.RWMutex{},
	}
}

func (i *InMemoryRecipeStore) GetRecipe(name string) *Recipe {
	i.lock.Lock()
	defer i.lock.Unlock()
	recipe, ok := i.recipesStore[name]
	if !ok {
		return nil
	}

	return &recipe
}

func (i *InMemoryRecipeStore) RecordRecipe(name string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.recipesStore[name] = Recipe{
		Title:       name,
		Description: "",
		Ingredients: []string{},
	}
}

func (i *InMemoryRecipeStore) GetUsers() []User {
	return i.usersStore
}

func (i *InMemoryRecipeStore) RecordUser(name string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.usersStore = append(i.usersStore, User{name})
}
