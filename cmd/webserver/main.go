package main

import (
	"racook/views/page"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

var counter int = 0

func main() {
	e := echo.New()
	e.Use()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, took=${latency_human}, error=${error}\n",
	}))

	e.Static("/static", "assets")

	e.GET("/", func(c echo.Context) error {
		return Render(c, 200, page.Home("0"))
	})

	e.POST("/couter", func(c echo.Context) error {
		counter++
		strCounter := strconv.Itoa(counter)

		return Render(c, 200, page.Counter(strCounter))
	})

	e.Logger.Fatal(e.Start(":6969"))
}
