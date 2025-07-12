package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var MRSizeHandler = OSSMetricHandler(data.BuildMRSizeQuery)
