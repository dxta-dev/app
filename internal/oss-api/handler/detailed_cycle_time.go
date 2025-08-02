package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var DetailedCycleTimeHandler = OSSMetricHandler(data.BuildDetailedCycleTimeQuery)
