package core

import (
	"context"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ContainerStats returns stats for a container stats request based on container id.
func (ds *templateService) ContainerStats(
	_ context.Context,
	r *runtimeapi.ContainerStatsRequest,
) (*runtimeapi.ContainerStatsResponse, error) {

	return &runtimeapi.ContainerStatsResponse{Stats: nil}, nil
}

// ListContainerStats returns stats for a list container stats request based on a filter.
func (ds *templateService) ListContainerStats(
	ctx context.Context,
	r *runtimeapi.ListContainerStatsRequest,
) (*runtimeapi.ListContainerStatsResponse, error) {

	return &runtimeapi.ListContainerStatsResponse{Stats: nil}, nil
}
