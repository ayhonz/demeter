package server

import (
	"fmt"
	"net/http"
	"racook/views/page"

	"github.com/labstack/echo/v4"
)

func (app *Application) HomePageHander(c echo.Context) error {

	recipe, err := app.Recipes.List()
	if err != nil {
		return err
	}

	return Render(c, 200, page.Home(recipe))
}

func (app *Application) CreateRecipePageHandler(c echo.Context) error {
	return Render(c, 200, page.CreateRecipe())
}

func (app *Application) CreateRecipeHandler(c echo.Context) error {
	var recipe RecipeForm
	err := c.Bind(&recipe)

	id, err := app.Recipes.Insert(recipe.Title, recipe.Description, recipe.Ingredients, recipe.Categories)
	if err != nil {
		return err
	}

	// Is this redirect correct? :?
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/recipes/%d", id))
	return c.String(http.StatusAccepted, "Created")
}

func (app *Application) GetDetailHandler(c echo.Context) error {
	id := c.Param("id")

	recipe, err := app.Recipes.Get(id)
	if err != nil {
		return err
	}

	return Render(c, 200, page.Detail(recipe))
}

func (app *Application) LoginPageHandler(c echo.Context) error {
	return Render(c, 200, page.Login())
}

func (app *Application) SignupPageHandler(c echo.Context) error {
	return Render(c, 200, page.Signup())
}

func (app *Application) SignupHandler(c echo.Context) error {
	return Render(c, 200, page.Signup())
}

func (app *Application) LoginHandler(c echo.Context) error {
	return Render(c, 200, page.Signup())
}

func (app *Application) LogoutHandler(c echo.Context) error {
	return Render(c, 200, page.Signup())
}
