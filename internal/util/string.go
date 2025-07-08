package util

import (
	"regexp"
	"strings"
)

// Sanitizes a string to conform to Turso database name requirements.
// https://docs.turso.tech/api-reference/databases/create#body-name
func SanitizeString(s string) string {
	s = strings.ToLower(s)

	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	reg = regexp.MustCompile(`^-+|-+$`)
	s = reg.ReplaceAllString(s, "")

	if len(s) > 64 {
		s = s[:64]
	}

	reg = regexp.MustCompile(`-+%`)
	s = reg.ReplaceAllString(s, "")

	return s
}
