package main

import (
	"net/http"

	"github.com/dxta-dev/app/internal/handler/api"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hell")
	})

	e.GET("/mr-size/:org/:repo", api.MrSizeHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
