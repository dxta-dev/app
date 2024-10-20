package host

import (
	"net/http"
	"strconv"
	"strings"
)

type GitHub struct {
	endpoint string
}

func NewGitHubHost() *GitHub {
	return &GitHub{
		endpoint: "https://api.github.com",
	}
}

func unwrapRatelimit(resp *http.Response) RateLimit {
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
			url:   url,
			value: page,
		}
	}

	return links
}

func (g GitHub) UnwrapResponse(resp *http.Response) (*UnwrappedResponse, error) {
	links := unwrapLink(resp)
	totalPages := links[Last].value
	if totalPages == 0 {
		totalPages = 1
	}
	return &UnwrappedResponse{
		Links: links,
		Pagination: Pagination{
			TotalPages: totalPages,
		},
	}, nil
}

func (g GitHub) UnwrapRequest(req *http.Request) (*UnwrappedRequest, error) {
	tenantId, err := unwrapTenantId(req)
	if err != nil {
		return nil, err
	}
	return &UnwrappedRequest{
		TenantId: tenantId,
	}, nil
}
