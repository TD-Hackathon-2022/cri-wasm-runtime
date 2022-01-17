package core

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// StopContainer stops a running container with a grace period (i.e., timeout).
func (ds *templateService) StopContainer(
	_ context.Context,
	r *v1.StopContainerRequest,
) (*v1.StopContainerResponse, error) {
	logrus.Infof("stop container: %s, container count : %d", r.GetContainerId(), len(ds.containerCache))
	defer logrus.Infof("end stop container, container: %s", r.GetContainerId())
	containerCache := ds.containerCache[r.GetContainerId()]
	if containerCache == nil {
		return nil, fmt.Errorf("cannot find container")
	}
	containerCache.status.State = v1.ContainerState_CONTAINER_EXITED
	containerCache.status.FinishedAt = ds.clock.Now().UnixNano()
	containerCache.status.ExitCode = 0
	return &v1.StopContainerResponse{}, nil
}
