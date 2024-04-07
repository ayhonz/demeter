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
