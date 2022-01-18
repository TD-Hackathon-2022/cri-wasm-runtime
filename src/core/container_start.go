package core

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// StartContainer starts the container.
func (ds *templateService) StartContainer(
	_ context.Context,
	r *v1.StartContainerRequest,
) (*v1.StartContainerResponse, error) {
	//defer logrus.Infof("end start container, container: %s", r.GetContainerId())
	containerCache := ds.containerCache[r.GetContainerId()]
	if containerCache == nil {
		logrus.Infof("cannot find container : %s", r.GetContainerId())
		return nil, fmt.Errorf("cannot find container")
	}

	logrus.Infof("start container: %s, container count : %d, containerName: %s, sandoxId: %s", r.GetContainerId(), len(ds.containerCache), containerCache.config.GetMetadata().GetName(), containerCache.sandboxId)
	containerCache.status.State = v1.ContainerState_CONTAINER_RUNNING
	containerCache.status.StartedAt = ds.clock.Now().UnixNano()
	return &v1.StartContainerResponse{}, nil
}
