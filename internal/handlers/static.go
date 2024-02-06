package handlers

import (
	"io/fs"
	"net/http"
	"dxta-dev/app"
	"github.com/labstack/echo/v4"
)


func (a App) PublicHandler() echo.HandlerFunc {
	publicFS, err := fs.Sub(static.Public, "public")

	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(publicFS))

	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000")

		fileServer.ServeHTTP(c.Response(), c.Request())

		return nil
	}
}
