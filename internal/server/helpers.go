package server

import (
	"context"
	"net/http"
	"racook/views"
	"time"
)

func (app *Application) newTemplateData(r *http.Request) views.TemplateData {
	return views.TemplateData{
		CurrentYear:   time.Now().Year(),
		Authenticated: app.isAuthenticated(r.Context()),
	}
}

func (app *Application) isAuthenticated(c context.Context) bool {
	return app.SessionManager.Exists(c, "authenticatedUserID")
}
