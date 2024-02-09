package handlers

import (
	"dxta-dev/app"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)


func (a App) PublicHandler() echo.HandlerFunc {
	publicFS, err := fs.Sub(static.Public, "public")

	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(publicFS))

	return func(c echo.Context) error {
		if !a.DebugMode {
			c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
		}

		fileServer.ServeHTTP(c.Response(), c.Request())

		return nil
	}
}
