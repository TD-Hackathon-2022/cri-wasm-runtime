package core

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// StopContainer stops a running container with a grace period (i.e., timeout).
func (ds *templateService) StopContainer(
	_ context.Context,
	r *v1.StopContainerRequest,
) (*v1.StopContainerResponse, error) {

	return &v1.StopContainerResponse{}, nil
}
