package server

import (
	"fmt"
	"log"
	"net/http"
	"racook/views"
	"racook/views/page"

	"github.com/a-h/templ"
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

	recipes, err := app.Recipes.List()
	if err != nil {
		return err
	}

	log.Println(c.Get("csrf"))
	templData := app.newTemplateData(c)
	templData.Recipes = recipes

	return Render(c, 200, page.Home(templData))
}

func (app *Application) CreateRecipePageHandler(c echo.Context) error {
	templData := app.newTemplateData(c)

	return Render(c, 200, page.CreateRecipe(templData))
}

func (app *Application) CreateRecipeHandler(c echo.Context) error {
	var recipe recipeForm
	err := c.Bind(&recipe)
	if err != nil {
		return err
	}
	// probably should move it somewhere else
	// but it will do for now :)
	userId := app.SessionManager.GetInt(c.Request().Context(), "authenticatedUserID")
	if userId == 0 {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	id, err := app.Recipes.Insert(userId, recipe.Title, recipe.Description, recipe.Ingredients, recipe.Categories)
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

	templData := app.newTemplateData(c)
	templData.Recipe = recipe

	return Render(c, 200, page.Detail(templData))
}

func (app *Application) LoginPageHandler(c echo.Context) error {
	templData := app.newTemplateData(c)

	return Render(c, 200, page.Login(templData))
}

func (app *Application) SignupPageHandler(c echo.Context) error {
	templData := app.newTemplateData(c)

	return Render(c, 200, page.Signup(templData))
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
	err := app.SessionManager.RenewToken(c.Request().Context())
	if err != nil {
		return err
	}

	app.SessionManager.Remove(c.Request().Context(), "authenticatedUserID")

	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(http.StatusOK, "Logged out")
}

func (app *Application) CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	var errorPage func(data views.TemplateData) templ.Component

	if code == http.StatusNotFound {
		errorPage = page.NotFound
	} else {
		errorPage = page.Error
	}

	templData := app.newTemplateData(c)
	Render(c, code, errorPage(templData))
}
