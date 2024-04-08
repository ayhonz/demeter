package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *Application) requireAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !app.isAuthenticated(c.Request().Context()) {
			c.Redirect(http.StatusSeeOther, "/user/login")
		}
		c.Response().Header().Add("Cache-control", "no-store")
		return next(c)
	}
}

func (app *Application) authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := app.SessionManager.GetInt(c.Request().Context(), "authenticatedUserID")
		if id == 0 {
			next(c)
			return nil
		}

		exists, err := app.Users.Exists(id)
		if err != nil {
			c.Error(err)
		}

		if exists {
		}
		return nil

	}
}
