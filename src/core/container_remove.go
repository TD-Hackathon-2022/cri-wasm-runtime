package core

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// RemoveContainer removes the container.
func (ds *templateService) RemoveContainer(
	_ context.Context,
	r *v1.RemoveContainerRequest,
) (*v1.RemoveContainerResponse, error) {
	// something here
	return &v1.RemoveContainerResponse{}, nil
}
