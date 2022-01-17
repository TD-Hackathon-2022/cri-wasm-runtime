package core

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ListContainers lists all containers matching the filter.
func (ds *templateService) ListContainers(
	_ context.Context,
	r *v1.ListContainersRequest,
) (*v1.ListContainersResponse, error) {
	// list all with filter
	// todo filter
	//logrus.Infof("list container, container count : %d", len(ds.containerCache))
	//defer logrus.Infof("end list container")
	items := make([]*v1.Container, 0, len(ds.containerCache))
	for containerId, containerCache := range ds.containerCache {
		item := &v1.Container{
			Id:           containerId,
			PodSandboxId: containerCache.sandboxId,
			Metadata:     containerCache.config.Metadata,
			Image:        ds.imageCache[containerCache.config.GetImage().GetImage()].image,
			ImageRef:     containerCache.config.GetImage().GetImage(),
			State:        containerCache.status.State,
			CreatedAt:    containerCache.status.CreatedAt,
			Labels:       containerCache.config.GetLabels(),
			Annotations:  containerCache.config.GetAnnotations(),
		}
		items = append(items, item)
	}
	return &v1.ListContainersResponse{Containers: items}, nil
}
