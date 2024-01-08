package static

import (
	"embed"
)

//go:embed public/*
var Public embed.FS

