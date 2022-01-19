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
	//defer logrus.Infof("end list container")
	filterSandboxId := r.GetFilter().GetPodSandboxId()
	filterContainerId := r.GetFilter().GetId()
	filterContainerState := r.GetFilter().GetState()
	items := make([]*v1.Container, 0, len(ds.containerCache))
	for containerId, containerCache := range ds.containerCache {
		var filterSuccess bool
		if len(filterSandboxId) != 0 {
			filterSuccess = filterSandboxId == containerCache.sandboxId
		}
		if len(filterContainerId) != 0 {
			filterSuccess = filterContainerId == containerId
		}
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
		if filterContainerState != nil {
			filterSuccess = filterContainerState.GetState() == containerCache.status.GetState()
			if filterSuccess {
				//logrus.Infof("list container: %s, container status : %d, image: %s", containerId, containerCache.status.GetState(), containerCache.config.GetImage().GetImage())
				items = append(items, item)
			}
		} else {
			items = append(items, item)
		}

	}
	//logrus.Infof("list container, sbid:{}, cid:{}, cs:{} : %s, %s ,%d, end list container, itemSize: %d", filterSandboxId, filterContainerId, filterContainerState.GetState(), len(items))
	return &v1.ListContainersResponse{Containers: items}, nil
}
