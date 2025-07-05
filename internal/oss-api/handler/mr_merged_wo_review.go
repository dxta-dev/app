package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var MRsMergedWithoutReviewHandler = OSSMetricHandler(data.BuildMRsMergedWithoutReviewQuery)
