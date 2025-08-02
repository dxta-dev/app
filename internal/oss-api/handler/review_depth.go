package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var ReviewDepthHandler = OSSMetricHandler(data.BuildReviewDepthQuery)
