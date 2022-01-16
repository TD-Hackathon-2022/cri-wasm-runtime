package core

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func (ds *templateService) UpdateContainerResources(
	_ context.Context,
	r *v1.UpdateContainerResourcesRequest,
) (*v1.UpdateContainerResourcesResponse, error) {

	return &v1.UpdateContainerResourcesResponse{}, nil
}
