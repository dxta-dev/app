package host

import (
	"reflect"
	"testing"
)

func TestUnwrapLink(t *testing.T) {
	tests := []struct {
		name       string
		linkHeader string
		expected   map[LinkKey]Link
	}{
		{
			name:       "Next and Last links",
			linkHeader: `<https://api.github.com/repositories/1300192/issues?page=4>; rel="next", <https://api.github.com/repositories/1300192/issues?page=597>; rel="last"`,
			expected: map[LinkKey]Link{
				Next: {
					url:   "https://api.github.com/repositories/1300192/issues?page=4",
					value: 4,
				},
				Last: {
					url:   "https://api.github.com/repositories/1300192/issues?page=597",
					value: 597,
				},
			},
		},
		{
			name:       "First, Prev, Next, and Last links",
			linkHeader: `<https://api.github.com/repositories/1300192/issues?page=2>; rel="prev", <https://api.github.com/repositories/1300192/issues?page=4>; rel="next", <https://api.github.com/repositories/1300192/issues?page=597>; rel="last", <https://api.github.com/repositories/1300192/issues?page=1>; rel="first"`,
			expected: map[LinkKey]Link{
				Previous: {
					url:   "https://api.github.com/repositories/1300192/issues?page=2",
					value: 2,
				},
				Next: {
					url:   "https://api.github.com/repositories/1300192/issues?page=4",
					value: 4,
				},
				Last: {
					url:   "https://api.github.com/repositories/1300192/issues?page=597",
					value: 597,
				},
				First: {
					url:   "https://api.github.com/repositories/1300192/issues?page=1",
					value: 1,
				},
			},
		},
		{
			name:       "Empty link header",
			linkHeader: ``,
			expected:   map[LinkKey]Link{},
		},
		{
			name:       "Gibberish link header",
			linkHeader: `<>; something="weird", <invalid-url?page=>; rel="next", random-text, <>; rel=invalid`,
			expected:   map[LinkKey]Link{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unwrapLink(tt.linkHeader)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unwrapLink() = %v, want %v", result, tt.expected)
			}
		})
	}
}
