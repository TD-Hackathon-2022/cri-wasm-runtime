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
	logrus.Infof("remove container: %s", r.GetContainerId())
	logrus.Infof("container count : %d", len(ds.containerCache))
	delete(ds.containerCache, r.GetContainerId())
	return &v1.RemoveContainerResponse{}, nil
}
