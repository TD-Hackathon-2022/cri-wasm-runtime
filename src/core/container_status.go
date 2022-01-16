package core

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ContainerStatus inspects the docker container and returns the status.
func (ds *templateService) ContainerStatus(
	_ context.Context,
	req *v1.ContainerStatusRequest,
) (*v1.ContainerStatusResponse, error) {
	//something here
	status := &v1.ContainerStatus{}
	return &v1.ContainerStatusResponse{Status: status}, nil
}
