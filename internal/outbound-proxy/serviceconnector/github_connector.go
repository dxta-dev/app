package host

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/dxta-dev/app/internal/assert"
)

type GitHubConnector struct {
	endpoint string
}

func NewGitHubConnector() *GitHubConnector {
	return &GitHubConnector{
		endpoint: "https://api.github.com",
	}
}

func unwrapRatelimit(resp *http.Response) RateLimit {
	assert.NotNil(resp, "Response must not be nil")

	resource := resp.Header.Get("X-Ratelimit-Resource")
	limit, _ := strconv.Atoi(resp.Header.Get("X-Ratelimit-Limit"))
	remaining, _ := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	used, _ := strconv.Atoi(resp.Header.Get("X-Ratelimit-Used"))
	reset, _ := strconv.ParseInt(resp.Header.Get("X-Ratelimit-Reset"), 10, 64)

	rateLimit := RateLimit{
		Resource:  resource,
		Limit:     limit,
		Remaining: remaining,
		RetryBy:   reset,
		Used:      used,
	}

	return rateLimit
}

func unwrapLink(resp *http.Response) map[LinkKey]Link {
	assert.NotNil(resp, "Request must not be nil")

	linkHeader := resp.Header.Get("link")
	links := make(map[LinkKey]Link)

	parts := strings.Split(linkHeader, ",")
	for _, part := range parts {
		sections := strings.Split(part, ";")
		if len(sections) < 2 {
			continue
		}
		url := strings.Trim(sections[0], " <>")

		var rel string
		for _, section := range sections[1:] {
			section = strings.TrimSpace(section)
			if strings.HasPrefix(section, "rel=") {
				rel = strings.Trim(section[4:], "\"")
				break
			}
		}

		urlParts := strings.Split(url, "page=")
		if len(urlParts) < 2 {
			continue
		}
		page, err := strconv.Atoi(urlParts[1])
		if err != nil {
			continue
		}

		linkKey, ok := stringToLinkKey(rel)
		if !ok {
			continue
		}

		links[linkKey] = Link{
			Url:   url,
			Value: page,
		}
	}

	return links
}

func (g GitHubConnector) UnwrapResponse(resp *http.Response) (*UnwrappedProxiedResponse, error) {
	links := unwrapLink(resp)
	totalPages := links[Last].Value
	if totalPages == 0 {
		totalPages = 1
	}
	rateLimit := unwrapRatelimit(resp)
	return &UnwrappedProxiedResponse{
		Links: links, Pagination: Pagination{
			TotalPages: totalPages,
		},
		RateLimit: rateLimit,
	}, nil
}

func (g GitHubConnector) UnwrapRequest(req *http.Request) *UnwrappedProxiedRequest {
	tenantId := unwrapTenantId(req)
	return &UnwrappedProxiedRequest{
		TenantId: tenantId,
	}
}

func (g GitHubConnector) MakeRequest(ctx context.Context, endpoint string, method string, headers http.Header, body []byte) (*http.Response, error) {
	return nil, nil
}
