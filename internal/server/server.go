package server

import (
	"net/http"
	"racook/internal/models"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	session "github.com/spazzymoto/echo-scs-session"
)

type Application struct {
	Recipes        *models.RecipeModel
	SessionManager *scs.SessionManager
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

	e.GET("/", app.HomePageHander, session.LoadAndSave(app.SessionManager))
	e.GET("/recipes/create", app.CreateRecipePageHandler, session.LoadAndSave(app.SessionManager))
	e.POST("/recipes", app.CreateRecipeHandler, session.LoadAndSave(app.SessionManager))
	e.GET("/recipes/:id", app.GetDetailHandler, session.LoadAndSave(app.SessionManager))

	e.GET("/login", app.LoginPageHandler, session.LoadAndSave(app.SessionManager))
	e.GET("/signup", app.SignupPageHandler, session.LoadAndSave(app.SessionManager))

	return e
}
