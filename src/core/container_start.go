package core

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// StartContainer starts the container.
func (ds *templateService) StartContainer(
	_ context.Context,
	r *v1.StartContainerRequest,
) (*v1.StartContainerResponse, error) {
	// something here
	// exec wasm
	return &v1.StartContainerResponse{}, nil
}
