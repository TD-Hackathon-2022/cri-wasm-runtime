package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/diannaowa/cri-template/streaming"
	"github.com/diannaowa/cri-template/utils"
	"google.golang.org/grpc/codes"
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

	timeout := time.Duration(utils.Min(req.Timeout, syncExecMaxTimeout)) * time.Second
	var stdoutBuffer, stderrBuffer bytes.Buffer
	err := ds.streamingRuntime.ExecWithContext(ctx, req.ContainerId, req.Cmd,
		nil, // in
		utils.WriteCloserWrapper(utils.LimitWriter(&stdoutBuffer, maxMsgSize)),
		utils.WriteCloserWrapper(utils.LimitWriter(&stderrBuffer, maxMsgSize)),
		false, // tty
		nil,   // resize
		timeout)

	// kubelet's backend runtime expects a grpc error with status code DeadlineExceeded on time out.
	if err == context.DeadlineExceeded {
		return nil, fmt.Errorf("deadline exceeded (%q): %v", codes.DeadlineExceeded, err.Error())
	}

	var exitCode int32
	if err != nil {
		exitError, ok := err.(utils.ExitError)
		if !ok {
			return nil, err
		}

		exitCode = int32(exitError.ExitStatus())
	}
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
