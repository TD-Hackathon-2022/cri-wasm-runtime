package core

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// CreateContainer creates a new container in the given PodSandbox
// Docker cannot store the log to an arbitrary location (yet), so we create an
// symlink at LogPath, linking to the actual path of the log.
func (ds *templateService) CreateContainer(
	_ context.Context,
	r *v1.CreateContainerRequest,
) (*v1.CreateContainerResponse, error) {
	//something here
	return &v1.CreateContainerResponse{
		ContainerId: "mock-uuid",
	}, nil
}
