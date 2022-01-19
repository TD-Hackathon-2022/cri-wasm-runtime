package core

import (
	"context"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// RunPodSandbox creates and starts a pod-level sandbox. Runtimes should ensure
// the sandbox is in ready state.
func (ds *templateService) RunPodSandbox(
	ctx context.Context,
	r *v1.RunPodSandboxRequest,
) (*v1.RunPodSandboxResponse, error) {
	//defer logrus.Infof("end run sandbox, sandboxId: %s", r.GetConfig().GetMetadata().GetUid())
	if r.GetConfig().GetLabels()["cri-runtime"] != "wasm" {
		return nil, fmt.Errorf("not wasm sandbox")
	}
	resp := &v1.RunPodSandboxResponse{PodSandboxId: uuid.New()}
	logrus.Infof("run sandbox, sandboxId: %s", resp.GetPodSandboxId())
	status := &v1.PodSandboxStatus{
		Id:        resp.GetPodSandboxId(),
		Metadata:  r.GetConfig().Metadata,
		State:     v1.PodSandboxState_SANDBOX_READY,
		CreatedAt: ds.clock.Now().UnixNano(),
		Network: &v1.PodSandboxNetworkStatus{
			Ip:            "172.31.19.45",
			AdditionalIps: make([]*v1.PodIP, 0),
		},
		Linux: &v1.LinuxPodSandboxStatus{
			Namespaces: &v1.Namespace{
				Options: r.GetConfig().GetLinux().GetSecurityContext().GetNamespaceOptions(),
			},
		},
		Labels:         r.GetConfig().Labels,
		Annotations:    r.GetConfig().Annotations,
		RuntimeHandler: r.RuntimeHandler,
	}
	ds.sandboxCache[resp.PodSandboxId] = &sandboxCacheModel{
		id:             resp.GetPodSandboxId(),
		config:         r.GetConfig(),
		status:         status,
		containerIdMap: make(map[string]string),
	}
	return resp, nil
}

/**
global cache:
podSandboxMap:{"PodSandboxId":PodSandbox}
1, RunPodSandbox --->PodSandboxId

containerMap:{"containerId":Container}

sandboxToContainer:{"sandboxid":["",""]}


trigger-->serverless-->scheduler--->call func(xx,xx,xx)


10W function-->function runtime -->scheduler alag...-->current node active jobs count

serverless
-->
*/
