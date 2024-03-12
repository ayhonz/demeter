package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewRecipeList(rdr io.Reader) ([]Recipe, error) {
	var recipes []Recipe
	err := json.NewDecoder(rdr).Decode(&recipes)
	if err != nil {
		err = fmt.Errorf("problem parsing recipes from file system, %v", err)
	}
	return recipes, err
}
