package core

import (
	"context"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// RemovePodSandbox removes the sandbox. If there are running containers in the
// sandbox, they should be forcibly removed.
func (ds *templateService) RemovePodSandbox(
	ctx context.Context,
	r *v1.RemovePodSandboxRequest,
) (*v1.RemovePodSandboxResponse, error) {
	logrus.Infof("remove sandbox, sandboxId: %s", r.GetPodSandboxId())
	logrus.Infof("sandbox count: %d", len(ds.sandboxCache))
	sandboxCache := ds.sandboxCache[r.GetPodSandboxId()]
	if sandboxCache != nil {
		delete(ds.sandboxCache, r.GetPodSandboxId())
	}
	return nil, nil
}
