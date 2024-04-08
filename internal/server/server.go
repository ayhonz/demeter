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
	Users          *models.UserModel
	SessionManager *scs.SessionManager
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
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:csrf",
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookiePath:     "/",
	}))

	e.Static("/static", "assets")

	e.Use(session.LoadAndSave(app.SessionManager))

	e.GET("/", app.HomePageHander)
	e.GET("/recipes/create", app.CreateRecipePageHandler, app.requireAuthentication)
	e.POST("/recipes", app.CreateRecipeHandler, app.requireAuthentication)
	e.GET("/recipes/:id", app.GetDetailHandler)

	e.GET("/user/login", app.LoginPageHandler)
	e.GET("/user/signup", app.SignupPageHandler)
	e.POST("/user/login", app.LoginHandler)
	e.POST("/user/signup", app.SignupHandler)
	e.POST("/user/logout", app.LogoutHandler, app.requireAuthentication)

	e.HTTPErrorHandler = app.CustomHTTPErrorHandler

	return e
}
