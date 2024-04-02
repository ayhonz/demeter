package server

import (
	"fmt"
	"net/http"
	"racook/internal/models"
	"racook/views/page"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	Recipes *models.RecipeModel
}

var counter int = 0

type RecipeForm struct {
	Title       string   `form:"title"`
	Description string   `form:"description"`
	Categories  []string `form:"categories[]"`
	Ingredients []string `form:"ingredients[]"`
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func (app *Application) Routes() http.Handler {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, took=${latency_human}, error=${error}\n",
	}))

	e.Static("/static", "assets")

	e.GET("/", func(c echo.Context) error {

		recipe, err := app.Recipes.List()
		if err != nil {
			return err
		}

		return Render(c, 200, page.Home(recipe))
	})

	e.GET("/recipes/create", func(c echo.Context) error {
		return Render(c, 200, page.CreateRecipe())
	})

	e.POST("/recipes", func(c echo.Context) error {
		var recipe RecipeForm
		err := c.Bind(&recipe)

		id, err := app.Recipes.Insert(recipe.Title, recipe.Description, recipe.Ingredients, recipe.Categories)
		if err != nil {
			return err
		}

		// Is this redirect correct? :?
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/recipes/%d", id))
		return c.String(http.StatusAccepted, "Created")
	})
	e.GET("/recipes/:id", func(c echo.Context) error {
		id := c.Param("id")

		recipe, err := app.Recipes.Get(id)
		if err != nil {
			return err
		}

		return Render(c, 200, page.Detail(recipe))
	})

	e.GET("/login", func(c echo.Context) error {
		return Render(c, 200, page.Login())
	})

	e.GET("/signup", func(c echo.Context) error {
		return Render(c, 200, page.Signup())
	})

	return e
}
