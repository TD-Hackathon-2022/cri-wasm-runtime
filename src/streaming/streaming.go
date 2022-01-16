/*
Copyright 2021 Mirantis

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package streaming

import (
	"context"
	"io"
	"time"

	"k8s.io/client-go/tools/remotecommand"

	"k8s.io/kubernetes/pkg/kubelet/cri/streaming"
)

type StreamingRuntime struct {
	ExecHandler ExecHandler
}

// ExecHandler knows how to execute a command in a running Docker container.
type ExecHandler interface {
	ExecInContainer(
		ctx context.Context,
		cmd []string,
		stdin io.Reader,
		stdout, stderr io.WriteCloser,
		tty bool,
		resize <-chan remotecommand.TerminalSize,
		timeout time.Duration,
	) error
}

var _ streaming.Runtime = &StreamingRuntime{}

func (r *StreamingRuntime) Exec(
	containerID string,
	cmd []string,
	in io.Reader,
	out, err io.WriteCloser,
	tty bool,
	resize <-chan remotecommand.TerminalSize,
) error {
	return r.ExecWithContext(context.TODO(), containerID, cmd, in, out, err, tty, resize, 0)
}

// ExecWithContext adds a context.
func (r *StreamingRuntime) ExecWithContext(
	ctx context.Context,
	containerID string,
	cmd []string,
	in io.Reader,
	out, errw io.WriteCloser,
	tty bool,
	resize <-chan remotecommand.TerminalSize,
	timeout time.Duration,
) error {

	return r.ExecHandler.ExecInContainer(
		ctx,
		cmd,
		in,
		out,
		errw,
		tty,
		resize,
		timeout,
	)
}

func (r *StreamingRuntime) Attach(
	containerID string,
	in io.Reader,
	out, errw io.WriteCloser,
	tty bool,
	resize <-chan remotecommand.TerminalSize,
) error {

	return nil
}

func (r *StreamingRuntime) PortForward(
	podSandboxID string,
	port int32,
	stream io.ReadWriteCloser,
) error {
	return nil
}
