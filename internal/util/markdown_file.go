package util

import (
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v2"
)

func ReadMarkdownFromFS[T any](fsys fs.FS, path string) (meta T, body string, err error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return meta, "", fmt.Errorf("failed to read %q: %w", path, err)
	}
	fm, b, err := splitFrontMatter(data)
	if err != nil {
		return meta, "", err
	}
	if err := yaml.Unmarshal(fm, &meta); err != nil {
		return meta, "", fmt.Errorf("parsing frontmatter into %T: %w", meta, err)
	}
	return meta, string(b), nil
}

func splitFrontMatter(data []byte) (fm []byte, body []byte, err error) {
	const delim = "---\n"
	text := string(data)
	if !strings.HasPrefix(text, delim) {
		return nil, data, nil
	}
	parts := strings.SplitN(text, "\n---\n", 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid frontmatter format")
	}
	fm = []byte(strings.TrimPrefix(parts[0], delim))
	body = []byte(parts[1])
	return fm, body, nil
}
