package handlers

import (
	"crypto/md5"
	"dxta-dev/app"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func PublicHandler() echo.HandlerFunc {
	publicFS, err := fs.Sub(static.Public, "public")

	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(publicFS))

	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000")

		path := c.Request().URL.Path

		path = "style.css"
		file, err := publicFS.Open(path)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		lastModified := fileInfo.ModTime().Format(time.RFC1123)

		etag :=  fmt.Sprintf("%x", md5.Sum([]byte(lastModified)))

		if match := c.Request().Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				c.Response().WriteHeader(http.StatusNotModified)
				return nil
			}
		}

		c.Response().Header().Set("Etag", etag)

		fileServer.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
