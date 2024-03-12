package handler

import (
	"github.com/donseba/go-htmx"
)

type App struct {
	HTMX           *htmx.HTMX
	BuildTimestamp string
	DebugMode      bool
}
