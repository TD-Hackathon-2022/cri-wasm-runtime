package core

import (
	"context"
	"io"
	"time"

	"k8s.io/client-go/tools/remotecommand"
)

// NativeExecHandler executes commands in Docker containers using Docker's exec API.
type NativeExecHandler struct{}

// ExecInContainer executes the cmd in container using the Docker's exec API
func (*NativeExecHandler) ExecInContainer(
	ctx context.Context,
	cmd []string,
	stdin io.Reader,
	stdout, stderr io.WriteCloser,
	tty bool,
	resize <-chan remotecommand.TerminalSize,
	timeout time.Duration,
) error {
	return nil
}
