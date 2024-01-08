package handlers

import (
	"io/fs"
	"net/http"
	"dxta-dev/app"
	"github.com/labstack/echo/v4"
)

func PublicHandler() echo.HandlerFunc {
	publicFS, err := fs.Sub(static.Public, "public")

	if err != nil {
		panic(err)
	}

	return echo.WrapHandler(http.FileServer(http.FS(publicFS)))
}
