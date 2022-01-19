package core

import (
	"context"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// RemoveContainer removes the container.
func (ds *templateService) RemoveContainer(
	_ context.Context,
	r *v1.RemoveContainerRequest,
) (*v1.RemoveContainerResponse, error) {
	// something here
	logrus.Infof("remove container: %s, container count : %d", r.GetContainerId(), len(ds.containerCache))
	//defer logrus.Infof("end remove container, container: %s", r.GetContainerId())
	containerCache := ds.containerCache[r.GetContainerId()]
	if containerCache != nil {
		delete(ds.containerCache, r.GetContainerId())
		delete(ds.sandboxCache[containerCache.sandboxId].containerIdMap, r.GetContainerId())
	}
	return &v1.RemoveContainerResponse{}, nil
}
