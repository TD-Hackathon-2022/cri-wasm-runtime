package core

import (
	"context"
	"github.com/diannaowa/cri-template/streaming"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// Attach prepares a streaming endpoint to attach to a running container, and returns the address.
func (ds *templateService) Attach(
	_ context.Context,
	req *v1.AttachRequest,
) (*v1.AttachResponse, error) {
	if ds.streamingServer == nil {
		return nil, streaming.NewErrorStreamingDisabled("attach")
	}

	return ds.streamingServer.GetAttach(req)
}
