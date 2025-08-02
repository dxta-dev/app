package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var CodingTimeHandler = OSSMetricHandler(data.BuildCodingTimeQuery)
