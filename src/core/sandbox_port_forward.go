package core

import (
	"context"
	"github.com/diannaowa/cri-template/streaming"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// PortForward prepares a streaming endpoint to forward ports from a PodSandbox, and returns the address.
func (ds *templateService) PortForward(
	_ context.Context,
	req *v1.PortForwardRequest,
) (*v1.PortForwardResponse, error) {
	if ds.streamingServer == nil {
		return nil, streaming.NewErrorStreamingDisabled("port forward")
	}
	return ds.streamingServer.GetPortForward(req)
}
