package server

import (
	"fmt"
	"net/http"
	"racook/internal/models"
	"racook/views/page"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	Recipes *models.RecipeModel
}

var counter int = 0

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
		strCounter := strconv.Itoa(counter)
		return Render(c, 200, page.Home(strCounter))
	})

	e.POST("/couter", func(c echo.Context) error {
		counter++
		strCounter := strconv.Itoa(counter)

		return Render(c, 200, page.Counter(strCounter))
	})
	e.POST("/recipes", func(c echo.Context) error {
		title := "test"
		description := "description"
		ingredients := []string{"ingredients"}
		categories := []string{"categories"}

		id, err := app.Recipes.Insert(title, description, ingredients, categories)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/recipes/%d", id))
	})
	e.GET("/recipes/:id", func(c echo.Context) error {
		id := c.Param("id")

		recipe, err := app.Recipes.Get(id)
		if err != nil {
			return err
		}

		return Render(c, 200, page.Detail(recipe))
	})

	return e
}
