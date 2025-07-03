package onboarding

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit/github_primary_ratelimit"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit/github_secondary_ratelimit"
	"github.com/google/go-github/v73/github"
)

type GithubConfig struct {
	GithubAppId         int64
	GithubAppPrivateKey []byte
	RoundTripper        http.RoundTripper
}

func LoadGithubConfig() (*GithubConfig, error) {
	appIdStr := os.Getenv("GITHUB_APP_ID")
	appPrivateKeyStr := os.Getenv("GITHUB_APP_PRIVATE_KEY")

	if appIdStr == "" {
		return nil, errors.New("GITHUB_APP_ID not set")
	}

	if appPrivateKeyStr == "" {
		return nil, errors.New("GITHUB_APP_PRIVATE_KEY not set")
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)

	if err != nil {
		return nil, errors.New("could not parse app id string to int64")
	}

	appPrivateKey, err := base64.StdEncoding.DecodeString(appPrivateKeyStr)

	if err != nil {
		return nil, errors.New("failed to decode base64 string")
	}

	return &GithubConfig{
		GithubAppId:         appId,
		GithubAppPrivateKey: appPrivateKey,
		RoundTripper:        http.DefaultTransport,
	}, nil
}

func getInstallationTransport(tr http.RoundTripper, installationId int64, appId int64, appPrivateKey []byte) (http.RoundTripper, error) {
	itt, err := ghinstallation.New(tr, appId, installationId, appPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("failed to create apps transport: %w", err)
	}

	return itt, nil
}

func getAppTransport(tr http.RoundTripper, appId int64, appPrivateKey []byte) (http.RoundTripper, error) {
	atr, err := ghinstallation.NewAppsTransport(tr, appId, appPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("failed to create apps transport: %w", err)
	}

	return atr, nil
}

func createLimiter(tr http.RoundTripper) http.RoundTripper {
	return github_ratelimit.New(tr,
		github_primary_ratelimit.WithLimitDetectedCallback(func(ctx *github_primary_ratelimit.CallbackContext) {
			fmt.Printf("Primary rate limit detected: category %s, reset time: %v\n", ctx.Category, ctx.ResetTime)
		}),
		github_secondary_ratelimit.WithLimitDetectedCallback(func(ctx *github_secondary_ratelimit.CallbackContext) {
			fmt.Printf("Secondary rate limit detected: reset time: %v, total sleep time: %v\n", ctx.ResetTime, ctx.TotalSleepTime)
		}),
	)
}

func NewInstallationClient(installationId int64, tr http.RoundTripper, cfg GithubConfig) (*github.Client, error) {
	tr, err := getInstallationTransport(tr, installationId, cfg.GithubAppId, cfg.GithubAppPrivateKey)

	if err != nil {
		return nil, err
	}

	tr = createLimiter(tr)

	return github.NewClient(&http.Client{Transport: tr}), nil
}

func InitAppClient(cfg GithubConfig) (*github.Client, error) {
	tr, err := getAppTransport(cfg.RoundTripper, cfg.GithubAppId, cfg.GithubAppPrivateKey)

	if err != nil {
		return nil, err
	}

	tr = createLimiter(tr)

	client := github.NewClient(&http.Client{Transport: tr})

	return client, nil
}
