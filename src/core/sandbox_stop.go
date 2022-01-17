package core

import (
	"context"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// StopPodSandbox stops the sandbox. If there are any running containers in the
// sandbox, they should be force terminated.
// better to cut our losses assuming an out of band GC routine will cleanup
// after us?
func (ds *templateService) StopPodSandbox(
	ctx context.Context,
	r *v1.StopPodSandboxRequest,
) (*v1.StopPodSandboxResponse, error) {
	logrus.Infof("stop sandbox, sandboxId: %s", r.GetPodSandboxId())
	logrus.Infof("sandbox count: %d", len(ds.sandboxCache))
	sandboxCache := ds.sandboxCache[r.GetPodSandboxId()]
	if sandboxCache == nil {
		return nil, nil
	}
	// todo status after stop
	sandboxCache.status.State = v1.PodSandboxState_SANDBOX_NOTREADY
	return nil, nil
}
