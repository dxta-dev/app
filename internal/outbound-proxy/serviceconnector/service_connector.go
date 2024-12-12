package host

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dxta-dev/app/internal/assert"
)

const (
	GITHUB int = iota
	GITLAB
	JIRA
	LINEAR
)

type ServiceConnector interface {
	UnwrapResponse(resp *http.Response) (*UnwrappedProxiedResponse, error)
	UnwrapRequest(req *http.Request) (*UnwrappedProxiedRequest, error)
	MakeRequest(ctx context.Context, endpoint string, method string, headers http.Header, body []byte) (*http.Response, error)
}

type LinkKey string

func stringToLinkKey(rel string) (LinkKey, bool) {
	links := map[string]LinkKey{
		"prev":  Previous,
		"next":  Next,
		"first": First,
		"last":  Last,
	}
	key, ok := links[rel]
	return key, ok
}

const (
	Previous LinkKey = "prev"
	Next     LinkKey = "next"
	First    LinkKey = "first"
	Last     LinkKey = "last"
)

type Pagination struct {
	TotalPages int
}

type UnwrappedProxiedResponse struct {
	Links map[LinkKey]Link
	Pagination
	RateLimit
}

type Link struct {
	Url   string
	Value int
}

type RateLimit struct {
	Limit     int
	Remaining int
	RetryBy   int64
	Resource  string
	Used      int
}

type UnwrappedProxiedRequest struct {
	TenantId string
}

func unwrapTenantId(req *http.Request) string {
	assert.NotNil(req, "Request must not be nil")

	tenantId := req.Header.Get("X-Tenant-Id")
	assert.Assert(tenantId != "", "Tenant ID must not be empty")

	return tenantId
}

func CreateResponse(unwrappedProxiedResponse *UnwrappedProxiedResponse, resp *http.Response) {
	_ = createResponseHeaders(unwrappedProxiedResponse)
}

func createResponseHeaders(unwrappedProxiedResponse *UnwrappedProxiedResponse) http.Header {
	assert.NotNil(unwrappedProxiedResponse, "Unwrapped proxied response must not be nil")

	headers := http.Header{}

	headers.Set("X-Pagination-Total", strconv.Itoa(unwrappedProxiedResponse.TotalPages))

	headers.Set("X-RateLimit-Limit", strconv.Itoa(unwrappedProxiedResponse.RateLimit.Limit))
	headers.Set("X-RateLimit-Remaining", strconv.Itoa(unwrappedProxiedResponse.RateLimit.Remaining))
	headers.Set("X-RateLimit-RetryBy", strconv.FormatInt(unwrappedProxiedResponse.RateLimit.RetryBy, 10))
	headers.Set("X-RateLimit-Resource", unwrappedProxiedResponse.RateLimit.Resource)
	headers.Set("X-RateLimit-Used", strconv.Itoa(unwrappedProxiedResponse.RateLimit.Used))

	for key, link := range unwrappedProxiedResponse.Links {
		headers.Set("X-Link-"+string(key), link.Url)
		headers.Set("X-Link-"+string(key)+"-Value", strconv.Itoa(link.Value))
	}

	return headers
}
