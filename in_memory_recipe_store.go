package main

import "sync"

func NewInMemoryRecipeStore() *InMemoryRecipeStore {
	return &InMemoryRecipeStore{
		map[string]string{},
		sync.RWMutex{},
	}
}

type InMemoryRecipeStore struct {
	store map[string]string
	lock  sync.RWMutex
}

func (i *InMemoryRecipeStore) GetRecipe(name string) string {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.store[name]
}
func (i *InMemoryRecipeStore) RecordRecipe(name string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.store[name] = name
}
