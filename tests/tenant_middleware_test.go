package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/utils"

	"github.com/labstack/echo/v4"
)

func TestTenantMiddleware(t *testing.T) {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	var mockDatabaseUrl = "libsql://john-cena"

	testCases := []struct {
		name              string
		host              string
		expectedIsRoot    bool
		expectedSubdomain string
	}{
		{"root_domain", "dxta.dev", true, "root"},
		{"subdomain", "crocoder.dxta.dev", false, "crocoder"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Host = tc.host
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			echoContext := e.NewContext(req, rec)

			var mockConfigTenants = make(map[string]utils.Tenant)
			if !tc.expectedIsRoot {
				mockConfigTenants[tc.expectedSubdomain] = utils.Tenant{
					Name:          tc.name,
					SubdomainName: tc.expectedSubdomain,
					DatabaseName:  tc.name,
					DatabaseUrl:   &mockDatabaseUrl,
				}
			}
			var mockConfig = utils.Config{
				IsMultiTenant:             true,
				ShouldUseSuperDatabase:    false,
				SuperDatabaseUrl:          nil,
				TenantDatabaseUrlTemplate: nil,
				Tenants:                   mockConfigTenants,
			}

			requestContext := echoContext.Request().Context()
			requestContext = middlewares.WithConfigContext(requestContext, &mockConfig)
			echoContext.SetRequest(echoContext.Request().WithContext(requestContext))

			if err := middlewares.MultiTenantMiddleware(func(c echo.Context) error { return nil })(echoContext); err != nil {
				t.Fatal(err)
			}

			requestContext = echoContext.Request().Context()

			isRoot, ok := requestContext.Value(middlewares.IsRootContext).(bool)
			if !ok {
				t.Errorf("is_root not set correctly")
			}
			if isRoot != tc.expectedIsRoot {
				t.Errorf("Expected is_root to be %v, got %v", tc.expectedIsRoot, isRoot)
			}

			subdomain, ok := requestContext.Value(middlewares.SubdomainContext).(string)
			if !ok {
				t.Errorf("subdomain not set correctly")
			}
			if subdomain != tc.expectedSubdomain {
				t.Errorf("Expected subdomain to be %v, got %v", tc.expectedSubdomain, subdomain)
			}

			_, ok = requestContext.Value(middlewares.TenantDatabaseURLContext).(string)
			if !ok && !tc.expectedIsRoot {
				t.Errorf("tenant database url not set correctly")
			}
			if ok && tc.expectedIsRoot {
				t.Errorf("tenant database url shouldn't be set for root")
			}
		})
	}
}
