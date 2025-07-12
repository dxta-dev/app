package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var ReviewTimeHandler = OSSMetricHandler(data.BuildReviewTimeQuery)
