package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dxta-dev/app/internal/web/template"

	"github.com/donseba/go-htmx"
)

type App struct {
	HTMX           *htmx.HTMX
	BuildTimestamp string
	DebugMode      bool
	Nonce          string
	State          State
}

type State struct {
	Team *int64
}

func (app *App) GenerateNonce() error {
	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return err
	}

	encodedNonce := hex.EncodeToString(nonce)
	app.Nonce = encodedNonce
	return nil
}

func (app *App) LoadState(r *http.Request) State {
	var team *int64
	if r.URL.Query().Has("team") {
		value, err := strconv.ParseInt(r.URL.Query().Get("team"), 10, 64)
		if err == nil {
			team = &value
		}
	}
	app.State.Team = team

	return app.State
}

func GetUrl(currentUrl string, params url.Values) (string, error) {
	if params == nil {
		params = url.Values{}
	}

	parsedURL, err := url.Parse(currentUrl)

	if err != nil {
		return "", err
	}

	requestUri := parsedURL.Path

	encodedParams := params.Encode()
	if encodedParams != "" {
		return fmt.Sprintf("%s?%s", requestUri, encodedParams), nil
	}

	return requestUri, nil
}

func (a *App) GetNavState() (template.NavState, error) {
	params := url.Values{}

	navState := template.GetDefaultNavState()

	if a.State.Team == nil {
		return navState, nil
	}

	if a.State.Team != nil {
		params.Add("team", fmt.Sprint(*a.State.Team))
	}

	rootUrl, err := GetUrl(navState.Root, params)
	if err != nil {
		return template.NavState{}, err
	}
	qMetricsUrl, err := GetUrl(navState.Metrics.Quality, params)
	if err != nil {
		return template.NavState{}, err
	}
	tMetricsUrl, err := GetUrl(navState.Metrics.Throughput, params)
	if err != nil {
		return template.NavState{}, err
	}

	return template.NavState{
		Root: rootUrl,
		Metrics: struct {
			Quality    string
			Throughput string
		}{
			Quality:    qMetricsUrl,
			Throughput: tMetricsUrl,
		},
	}, nil

}

func (app *App) GetUrlAppState(currentUrl string, params url.Values) (string, error) {
	if params == nil {
		params = url.Values{}
	}

	if app.State.Team != nil && !params.Has("team") {
		params.Add("team", fmt.Sprint(*app.State.Team))
	}

	return GetUrl(currentUrl, params)
}
