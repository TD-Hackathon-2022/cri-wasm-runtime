package core

import (
	"context"
	"github.com/sirupsen/logrus"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ImageFsInfo returns information of the filesystem that is used to store images.
func (ds *templateService) ImageFsInfo(
	_ context.Context,
	_ *runtimeapi.ImageFsInfoRequest,
) (*runtimeapi.ImageFsInfoResponse, error) {
	logrus.Infof("ImageFsInfo")
	return &runtimeapi.ImageFsInfoResponse{
		ImageFilesystems: nil,
	}, nil
}
