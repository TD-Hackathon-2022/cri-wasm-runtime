package core

import (
	"context"
	"fmt"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// PodSandboxStatus returns the status of the PodSandbox.
func (ds *templateService) PodSandboxStatus(
	ctx context.Context,
	req *v1.PodSandboxStatusRequest,
) (*v1.PodSandboxStatusResponse, error) {
	sandboxCache := ds.sandboxCache[req.GetPodSandboxId()]
	if sandboxCache == nil {
		return nil, fmt.Errorf("cannot find pod sandbox")
	}
	status := ds.sandboxCache[req.GetPodSandboxId()].status
	return &v1.PodSandboxStatusResponse{Status: status}, nil
}
