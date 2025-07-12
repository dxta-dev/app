package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var CycleTimeHandler = OSSMetricHandler(data.BuildCycleTimeQuery)
