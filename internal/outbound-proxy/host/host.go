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
	switch rel {
	case "prev":
		return Previous, true
	case "next":
		return Next, true
	case "first":
		return First, true
	case "last":
		return Last, true
	default:
		return "", false
	}
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
}

type Link struct {
	url   string
	value int
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
