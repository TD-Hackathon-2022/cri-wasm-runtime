package core

import (
	"context"
	"fmt"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ContainerStatus inspects the docker container and returns the status.
func (ds *templateService) ContainerStatus(
	_ context.Context,
	req *v1.ContainerStatusRequest,
) (*v1.ContainerStatusResponse, error) {
	//something here
	containerCache := ds.containerCache[req.GetContainerId()]
	//logrus.Infof("status container: %s, container status : %d, image: %s", req.GetContainerId(), containerCache.status.GetState(), containerCache.config.GetImage().GetImage())
	if containerCache == nil {
		return nil, fmt.Errorf("cannot find container")
	}
	status := containerCache.status
	return &v1.ContainerStatusResponse{Status: status}, nil
}
