package docs

import "embed"

//go:embed */*.md
var DocsFS embed.FS

type MarkdownHandlerFrontmatter struct {
	Title string `yaml:"title"`
}
