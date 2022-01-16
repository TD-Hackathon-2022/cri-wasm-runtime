// +build !linux,!windows

package core

import (
	"context"
	"fmt"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ImageFsInfo returns information of the filesystem that is used to store images.
func (ds *templateService) ImageFsInfo(
	_ context.Context,
	r *runtimeapi.ImageFsInfoRequest,
) (*runtimeapi.ImageFsInfoResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
