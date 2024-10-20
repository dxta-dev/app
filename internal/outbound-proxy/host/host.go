package host

import (
	"errors"
	"net/http"
)

const (
	GITHUB int = iota
	GITLAB
	JIRA
	LINEAR
)

type Host interface {
	UnwrapResponse(resp *http.Response) (*UnwrappedResponse, error)
	UnwrapRequest(req *http.Request) (*UnwrappedRequest, error)
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
