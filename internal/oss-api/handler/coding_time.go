package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var CodeTimeHandler = OSSMetricHandler(data.BuildCodingTimeQuery)
