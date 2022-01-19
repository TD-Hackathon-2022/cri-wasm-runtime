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
	//defer logrus.Infof("end stop sandbox, sandboxId: %s", r.GetPodSandboxId())
	sandboxCache := ds.sandboxCache[r.GetPodSandboxId()]
	if sandboxCache == nil {
		return nil, nil
	}
	sandboxCache.status.State = v1.PodSandboxState_SANDBOX_NOTREADY
	for containerId, _ := range sandboxCache.containerIdMap {
		ds.containerCache[containerId].status.State = v1.ContainerState_CONTAINER_EXITED
		ds.containerCache[containerId].status.FinishedAt = ds.clock.Now().UnixNano()
		ds.containerCache[containerId].status.ExitCode = 0
		ds.containerCache[containerId].status.Reason = "Completed"
	}
	return &v1.StopPodSandboxResponse{}, nil
}
