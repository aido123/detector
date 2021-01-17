package detectorapi

import (
	"context"

	"github.com/aido123/detector/pkg/containerservice/detector"
)

type DetectorClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, resourceName string, detectId string, startTime string, endTime string) (result detector.Detector, err error)
}

var _ DetectorClientAPI = (*detector.DetectorClient)(nil)
