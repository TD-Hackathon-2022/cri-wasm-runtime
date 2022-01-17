package core

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ListPodSandbox returns a list of Sandbox.
func (ds *templateService) ListPodSandbox(
	_ context.Context,
	r *v1.ListPodSandboxRequest,
) (*v1.ListPodSandboxResponse, error) {
	// todo filter
	//logrus.Infof("list sandbox, sandbox count: %d", len(ds.sandboxCache))
	filterSandboxId := r.GetFilter().GetId()
	filterSandboxState := r.GetFilter().GetState()
	items := make([]*v1.PodSandbox, 0, len(ds.sandboxCache))
	for id, cache := range ds.sandboxCache {
		filterSuccess := filterSandboxId == id && filterSandboxState.GetState() == cache.status.GetState()
		if filterSuccess {
			item := &v1.PodSandbox{
				Id:             id,
				Metadata:       cache.config.Metadata,
				State:          cache.status.State,
				CreatedAt:      cache.status.CreatedAt,
				Labels:         cache.config.Labels,
				Annotations:    cache.config.Annotations,
				RuntimeHandler: cache.status.GetRuntimeHandler(),
			}
			items = append(items, item)
		}
	}
	//logrus.Infof("end list2 sandbox")
	//logrus.Infof("end list sandbox, itemSize: %d", len(items))
	return &v1.ListPodSandboxResponse{Items: items}, nil
}
