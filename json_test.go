package cookbook

import (
	// "net/http/httptest"
	"testing"
)

func TestJsonResponse(t *testing.T) {
	// t.Run("returns a json response with 200 status", func(t *testing.T) {
	// 	response := httptest.NewRecorder()
	//
	// 	recipe := Recipe{
	// 		Title:       "chicken recipe",
	// 		Description: "",
	// 		Ingredients: []string{},
	// 	}
	//
	// 	responseWithJSON(response, 200, recipe)
	//
	// 	assertStatus(t, response.Code, 200)
	// 	assertContentType(t, response, jsonContentType)
	// 	assertRecipe(t, getRecipeFromResponse(t, response.Body), recipe)
	// })
	//
	// t.Run("error returns in correct format", func(t *testing.T) {
	// 	response := httptest.NewRecorder()
	//
	// 	responseWithError(response, 400, "invalid request")
	//
	// 	assertStatus(t, response.Code, 400)
	// 	assertContentType(t, response, jsonContentType)
	// 	assertResponseBody(t, response.Body.String(), "{\"error\":\"invalid request\"}")
	// })
}
