package core

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ContainerStatus inspects the docker container and returns the status.
func (ds *templateService) ContainerStatus(
	_ context.Context,
	req *v1.ContainerStatusRequest,
) (*v1.ContainerStatusResponse, error) {
	//something here
	logrus.Infof("status container: %s, container count : %d", req.GetContainerId(), len(ds.containerCache))
	containerCache := ds.containerCache[req.GetContainerId()]
	if containerCache.config.GetMetadata().GetName() == "nginx" {
		logrus.Infof("status container: %s, container status : %d", req.GetContainerId(), containerCache.status.GetState())
	}
	if containerCache == nil {
		return nil, fmt.Errorf("cannot find container")
	}
	status := containerCache.status
	return &v1.ContainerStatusResponse{Status: status}, nil
}
