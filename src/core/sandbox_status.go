package core

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// PodSandboxStatus returns the status of the PodSandbox.
func (ds *templateService) PodSandboxStatus(
	ctx context.Context,
	req *v1.PodSandboxStatusRequest,
) (*v1.PodSandboxStatusResponse, error) {
	logrus.Infof("status sandbox, sandboxId: %s", req.GetPodSandboxId())
	logrus.Infof("sandbox count: %d", len(ds.sandboxCache))
	sandboxCache := ds.sandboxCache[req.GetPodSandboxId()]
	if sandboxCache == nil {
		return nil, fmt.Errorf("cannot find pod sandbox")
	}
	status := ds.sandboxCache[req.GetPodSandboxId()].status
	return &v1.PodSandboxStatusResponse{Status: status}, nil
}
