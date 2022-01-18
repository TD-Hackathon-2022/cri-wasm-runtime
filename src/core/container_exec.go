package core

import (
	"bytes"
	"context"
	"github.com/diannaowa/cri-template/streaming"
	"github.com/diannaowa/cri-template/utils"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"time"
)

// ExecSync executes a command in the container, and returns the stdout output.
// If command exits with a non-zero exit code, an error is returned.
const syncExecMaxTimeout = 10

func (ds *templateService) ExecSync(
	ctx context.Context,
	req *v1.ExecSyncRequest,
) (*v1.ExecSyncResponse, error) {
	//logrus.Infof("exec container, container: %s, container count : %d", req.GetContainerId(), len(ds.containerCache))
	timeout := time.Duration(utils.Min(req.Timeout, syncExecMaxTimeout)) * time.Second
	logrus.Info(timeout)

	var stdoutBuffer, stderrBuffer bytes.Buffer
	var exitCode int32 = 0
	stdoutBuffer.WriteString("hello world")
	stderrBuffer.WriteString("")

	return &v1.ExecSyncResponse{
		Stdout:   stdoutBuffer.Bytes(),
		Stderr:   stderrBuffer.Bytes(),
		ExitCode: exitCode,
	}, nil
}

// Exec prepares a streaming endpoint to execute a command in the container, and returns the address.
func (ds *templateService) Exec(
	_ context.Context,
	req *v1.ExecRequest,
) (*v1.ExecResponse, error) {
	if ds.streamingServer == nil {
		return nil, streaming.NewErrorStreamingDisabled("exec")
	}
	return ds.streamingServer.GetExec(req)
}
