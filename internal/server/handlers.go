package server

import (
	"fmt"
	"net/http"
	"racook/views/page"

	"github.com/labstack/echo/v4"
)

type recipeForm struct {
	Title       string   `form:"title"`
	Description string   `form:"description"`
	Categories  []string `form:"categories[]"`
	Ingredients []string `form:"ingredients[]"`
}

type userSigunForm struct {
	Email     string `form:"email"`
	Password  string `form:"password"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
}

type userLoginform struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

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
	var recipe recipeForm
	err := c.Bind(&recipe)
	if err != nil {
		return err
	}

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
	var user userSigunForm
	err := c.Bind(&user)
	if err != nil {
		return err
	}
	err = app.Users.Insert(user.Email, user.FirstName, user.LastName, user.Password)
	if err != nil {
		return err
	}

	// prob want to do it on client with HX
	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(http.StatusAccepted, "Created")
}

func (app *Application) LoginHandler(c echo.Context) error {
	var user userLoginform
	err := c.Bind(&user)
	if err != nil {
		return err
	}

	id, err := app.Users.Authenticate(user.Email, user.Password)
	if err != nil {
		return err
	}

	err = app.SessionManager.RenewToken(c.Request().Context())
	if err != nil {
		return err
	}
	app.SessionManager.Put(c.Request().Context(), "authenticatedUserID", id)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(http.StatusOK, "Logged in")
}

func (app *Application) LogoutHandler(c echo.Context) error {
	return Render(c, 200, page.Signup())
}
