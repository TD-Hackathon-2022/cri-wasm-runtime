package core

import (
	"context"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// CreateContainer creates a new container in the given PodSandbox
// Docker cannot store the log to an arbitrary location (yet), so we create an
// symlink at LogPath, linking to the actual path of the log.
func (ds *templateService) CreateContainer(
	_ context.Context,
	r *v1.CreateContainerRequest,
) (*v1.CreateContainerResponse, error) {
	//something here
	//logrus.Infof("create container, container count : %d", len(ds.containerCache))
	//logrus.Infof("create container, image : %s", r.GetConfig().GetImage().GetImage())
	//defer logrus.Infof("end create container")
	containerId := uuid.New()
	sandboxId := r.GetPodSandboxId()
	containerConfig := r.GetConfig()
	containerStatus := &v1.ContainerStatus{
		Id:          containerId,
		Metadata:    r.GetConfig().GetMetadata(),
		State:       v1.ContainerState_CONTAINER_CREATED,
		CreatedAt:   ds.clock.Now().UnixNano(),
		Image:       ds.imageCache[r.GetConfig().GetImage().GetImage()].image,
		ImageRef:    r.GetConfig().GetImage().GetImage(),
		Labels:      r.GetConfig().GetLabels(),
		Annotations: r.GetConfig().GetAnnotations(),
	}

	sanboxCache := ds.sandboxCache[sandboxId]
	for sContainerId, _ := range sanboxCache.containerIdMap {
		if ds.containerCache[sContainerId] != nil && ds.containerCache[sContainerId].config.GetImage().GetImage() == r.GetConfig().GetImage().GetImage() {
			logrus.Infof("repeat create container, sandbox:%s image: %s", sandboxId, r.GetConfig().GetImage().GetImage())
			return &v1.CreateContainerResponse{
				ContainerId: sContainerId,
			}, nil
		}
	}
	containerCache := &containerCacheModel{
		id:            containerId,
		config:        containerConfig,
		status:        containerStatus,
		sandboxId:     sandboxId,
		sandboxConfig: ds.sandboxCache[sandboxId].config,
	}

	ds.containerCache[containerId] = containerCache
	ds.sandboxCache[sandboxId].containerIdMap[containerId] = containerId

	return &v1.CreateContainerResponse{
		ContainerId: containerId,
	}, nil
}
