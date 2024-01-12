package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dxta-dev/app/internal/middlewares"

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

			echoContext := e.NewContext(req, rec)
			var dummyTenantsMap = make(middlewares.TenantMap)
			if !tc.expectedIsRoot {
				dummyTenantsMap[tc.expectedSubdomain] = true
			}
			middlewares.LoadTenantsDummy(dummyTenantsMap)

			if err := middlewares.TenantMiddleware(func(c echo.Context) error { return nil })(echoContext); err != nil {
				t.Fatal(err)
			}

			context := echoContext.Request().Context()

			isRoot, ok := context.Value(middlewares.IsRootContext).(bool)
			if !ok {
				t.Errorf("is_root not set correctly")
			}
			if isRoot != tc.expectedIsRoot {
				t.Errorf("Expected is_root to be %v, got %v", tc.expectedIsRoot, isRoot)
			}

			subdomain, ok := context.Value(middlewares.SubdomainContext).(string)
			if !ok {
				t.Errorf("subdomain not set correctly")
			}
			if subdomain != tc.expectedSubdomain {
				t.Errorf("Expected subdomain to be %v, got %v", tc.expectedSubdomain, subdomain)
			}

			tenant, ok := context.Value(middlewares.TenantContext).(string)
			if !ok && !tc.expectedIsRoot {
				t.Errorf("tenant not set correctly")
			}
			if ok && tc.expectedIsRoot {
				t.Errorf("tenant root shouldn't be set")
			}
			if !tc.expectedIsRoot && tenant != tc.expectedSubdomain {
				t.Errorf("Expected tenant to be %v, got %v", tc.expectedSubdomain, tenant)
			}
		})
	}
}
