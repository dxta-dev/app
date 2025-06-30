package activities

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
	"github.com/google/go-github/v72/github"
)

var GithubConfig *GithubCfg

type GithubCfg struct {
	GithubAppId         int64
	GithubAppPrivateKey []byte
	GithubAppClient     *github.Client
	RoundTripper        http.RoundTripper
}

func LoadGithubConfig() (*GithubCfg, error) {
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

	GithubConfig = &GithubCfg{
		GithubAppId:         appId,
		GithubAppPrivateKey: appPrivateKey,
		GithubAppClient:     nil,
		RoundTripper:        http.DefaultTransport,
	}

	return GithubConfig, nil
}

func getInstallationTransport(tr http.RoundTripper, installationId int64) (http.RoundTripper, error) {
	itt, err := ghinstallation.New(tr, GithubConfig.GithubAppId, installationId, GithubConfig.GithubAppPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("failed to create apps transport: %w", err)
	}

	return itt, nil
}

func getAppTransport(tr http.RoundTripper) (http.RoundTripper, error) {
	atr, err := ghinstallation.NewAppsTransport(tr, GithubConfig.GithubAppId, GithubConfig.GithubAppPrivateKey)

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

func NewInstallationClient(installationId int64) (*github.Client, error) {
	tr := GithubConfig.RoundTripper
	tr, err := getInstallationTransport(tr, installationId)

	if err != nil {
		return nil, err
	}

	tr = createLimiter(tr)

	return github.NewClient(&http.Client{Transport: tr}), nil
}

func InitAppClient() error {
	tr := GithubConfig.RoundTripper
	tr, err := getAppTransport(tr)

	if err != nil {
		return err
	}

	tr = createLimiter(tr)

	GithubConfig.GithubAppClient = github.NewClient(&http.Client{Transport: tr})

	return nil
}

type NewInstallationClientFunc func(installationId int64) (*github.Client, error)
type GithubActivities struct {
	githubConfig          GithubCfg
	newInstallationClient NewInstallationClientFunc
}

func InitGHActivities(githubCfg GithubCfg) *GithubActivities {
	return &GithubActivities{githubConfig: githubCfg, newInstallationClient: NewInstallationClient}
}
