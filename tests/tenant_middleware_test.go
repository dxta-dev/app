package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dxta-dev/app/internals/middlewares"

	"github.com/labstack/echo/v4"
)

func TestTenantMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(middlewares.TenantMiddleware)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	testCases := []struct {
		name              string
		host              string
		expectedIsRoot    bool
		expectedSubdomain string
	}{
		{"RootDomain", "dxta.dev", true, "root"},
		{"SubDomain", "crocoder.dxta.dev", false, "crocoder"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Host = tc.host
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			context := e.NewContext(req, rec)

			if err := middlewares.TenantMiddleware(func(c echo.Context) error { return nil })(context); err != nil {
				t.Fatal(err)
			}

			isRoot, ok := context.Get("is_root").(bool)
			if !ok {
				t.Errorf("is_root not set correctly")
			}
			if isRoot != tc.expectedIsRoot {
				t.Errorf("Expected is_root to be %v, got %v", tc.expectedIsRoot, isRoot)
			}

			subdomain, ok := context.Get("subdomain").(string)
			if !ok {
				t.Errorf("subdomain not set correctly")
			}
			if subdomain != tc.expectedSubdomain {
				t.Errorf("Expected subdomain to be %v, got %v", tc.expectedSubdomain, subdomain)
			}
		})
	}
}
