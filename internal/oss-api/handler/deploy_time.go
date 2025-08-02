package handler

import (
	"github.com/dxta-dev/app/internal/data"
)

var DeployTimeHandler = OSSMetricHandler(data.BuildDeployTimeQuery)
