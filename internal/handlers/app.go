package handlers

import (
	"dxta-dev/app/internal/utils"

	"github.com/donseba/go-htmx"
)

type App struct {
	HTMX   *htmx.HTMX
	Config *utils.Config
}
