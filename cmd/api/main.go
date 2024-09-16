package main

import (
	"net/http"

	"github.com/dxta-dev/app/internal/handler/api"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hell")
	})

	e.GET("/code-chage/:org/:repo", api.CodeChangeHandler)
	e.GET("/coding-time/:org/:repo", api.CodingTimeHandler)
	e.GET("/commits/:org/:repo", api.CommitsHandler)
	e.GET("/cycle-time/:org/:repo", api.CycleTimeHandler)
	e.GET("/deploy-freq/:org/:repo", api.DeployFrequencyHandler)
	e.GET("/deploy-time/:org/:repo", api.DeployTimeHandler)
	e.GET("/handover/:org/:repo", api.HandoverHandler)
	e.GET("/merge-freq/:org/:repo", api.MergeFrequencyHandler)
	e.GET("/mr-merged-wo-review/:org/:repo", api.MRsMergedWithoutReviewHandler)
	e.GET("/mr-opened/:org/:repo", api.MRSOpenedHandler)
	e.GET("/mr-pickup-time/:org/:repo", api.MRPickupTimeHandler)
	e.GET("/mr-size/:org/:repo", api.MRSizeHandler)
	e.GET("/review/:org/:repo", api.ReviewHandler)
	e.GET("/review-depth/:org/:repo", api.ReviewDepthHandler)
	e.GET("/review-time/:org/:repo", api.ReviewTimeHandler)
	e.GET("/time-to-merge/:org/:repo", api.TimeToMergeHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
