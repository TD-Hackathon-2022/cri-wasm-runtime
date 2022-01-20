package core

import (
	"context"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ImageFsInfo returns information of the filesystem that is used to store images.
func (ds *templateService) ImageFsInfo(
	_ context.Context,
	_ *runtimeapi.ImageFsInfoRequest,
) (*runtimeapi.ImageFsInfoResponse, error) {
	return &runtimeapi.ImageFsInfoResponse{
		ImageFilesystems: nil,
	}, nil
}
