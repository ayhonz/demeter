package server

import (
	"context"
	"racook/views"
	"time"

	"github.com/labstack/echo/v4"
)

func (app *Application) newTemplateData(c echo.Context) views.TemplateData {
	return views.TemplateData{
		CurrentYear:   time.Now().Year(),
		Authenticated: app.isAuthenticated(c.Request().Context()),
		CRSFToken:     c.Get("csrf").(string),
	}
}

func (app *Application) isAuthenticated(c context.Context) bool {
	return app.SessionManager.Exists(c, "authenticatedUserID")
}
