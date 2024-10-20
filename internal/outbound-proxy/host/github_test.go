package host

import (
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
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
			resp := &http.Response{
				Header: http.Header{
					"Link": []string{tt.linkHeader},
				},
			}
			result := unwrapLink(resp)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unwrapLink() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnwrapRatelimit(t *testing.T) {
	resetTime := time.Date(2022, 3, 18, 0, 0, 0, 0, time.UTC).Unix()

	tests := []struct {
		name     string
		headers  map[string]string
		expected RateLimit
	}{
		{
			name: "Standard Rate Limit Headers",
			headers: map[string]string{
				"X-Ratelimit-Resource":  "core",
				"X-Ratelimit-Limit":     "5000",
				"X-Ratelimit-Remaining": "4995",
				"X-Ratelimit-Used":      "5",
				"X-Ratelimit-Reset":     strconv.FormatInt(resetTime, 10),
			},
			expected: RateLimit{
				Resource:  "core",
				Limit:     5000,
				Remaining: 4995,
				Used:      5,
				RetryBy:   resetTime,
			},
		},
		{
			name: "Rate Limit with zero remaining",
			headers: map[string]string{
				"X-Ratelimit-Resource":  "search",
				"X-Ratelimit-Limit":     "100",
				"X-Ratelimit-Remaining": "0",
				"X-Ratelimit-Used":      "100",
				"X-Ratelimit-Reset":     strconv.FormatInt(resetTime, 10),
			},
			expected: RateLimit{
				Resource:  "search",
				Limit:     100,
				Remaining: 0,
				Used:      100,
				RetryBy:   resetTime,
			},
		},
		{
			name:    "Empty Headers",
			headers: map[string]string{},
			expected: RateLimit{
				Resource:  "",
				Limit:     0,
				Remaining: 0,
				Used:      0,
				RetryBy:   0,
			},
		},
		{
			name: "Malformed Header Values",
			headers: map[string]string{
				"X-Ratelimit-Resource":  "graphql",
				"X-Ratelimit-Limit":     "not-a-number",
				"X-Ratelimit-Remaining": "abc",
				"X-Ratelimit-Used":      "def",
				"X-Ratelimit-Reset":     "not-a-timestamp",
			},
			expected: RateLimit{
				Resource:  "graphql",
				Limit:     0,
				Remaining: 0,
				Used:      0,
				RetryBy:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				Header: make(http.Header),
			}
			for key, value := range tt.headers {
				resp.Header.Set(key, value)
			}

			result := unwrapRatelimit(resp)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unwrapRatelimit() = %v, want %v", result, tt.expected)
			}
		})
	}
}
