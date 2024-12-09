package host

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

const (
	GITHUB int = iota
	GITLAB
	JIRA
	LINEAR
)

type ServiceConnector interface {
	UnwrapResponse(resp *http.Response) (*UnwrappedResponse, error)
	UnwrapRequest(req *http.Request) (*UnwrappedRequest, error)
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

type UnwrappedResponse struct {
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

type UnwrappedRequest struct {
	TenantId string
}

func unwrapTenantId(req *http.Request) (string, error) {
	tenantId := req.Header.Get("X-Tenant-Id")
	if tenantId == "" {
		return "", errors.New("")
	}
	return tenantId, nil
}

func CreateResponse(host ServiceConnector, resp *http.Response) (any, error) {
	_, err := CreateResponseHeaders(host, resp)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func CreateResponseHeaders(serviceConnector ServiceConnector, resp *http.Response) (http.Header, error) {
	unwrappedResponse, err := serviceConnector.UnwrapResponse(resp)
	if err != nil {
		return nil, err
	}
	headers := http.Header{}

	headers.Set("X-Pagination-Total", strconv.Itoa(unwrappedResponse.TotalPages))

	headers.Set("X-RateLimit-Limit", strconv.Itoa(unwrappedResponse.RateLimit.Limit))
	headers.Set("X-RateLimit-Remaining", strconv.Itoa(unwrappedResponse.RateLimit.Remaining))
	headers.Set("X-RateLimit-RetryBy", strconv.FormatInt(unwrappedResponse.RateLimit.RetryBy, 10))
	headers.Set("X-RateLimit-Resource", unwrappedResponse.RateLimit.Resource)
	headers.Set("X-RateLimit-Used", strconv.Itoa(unwrappedResponse.RateLimit.Used))

	for key, link := range unwrappedResponse.Links {
		headers.Set("X-Link-"+string(key), link.Url)
		headers.Set("X-Link-"+string(key)+"-Value", strconv.Itoa(link.Value))
	}

	return headers, nil
}
