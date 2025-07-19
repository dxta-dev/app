package handler

import (
	"github.com/dxta-dev/app/internal/data"
)


var HandoverHandler = OSSMetricHandler(data.BuildHandoverQuery)
