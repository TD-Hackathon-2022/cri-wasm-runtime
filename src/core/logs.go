package core

import (
	"context"
	"fmt"
	"github.com/diannaowa/cri-template/config"
	"io"
	v1 "k8s.io/api/core/v1"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ReopenContainerLog reopens the container log file.
func (ds *templateService) ReopenContainerLog(
	_ context.Context,
	_ *runtimeapi.ReopenContainerLogRequest,
) (*runtimeapi.ReopenContainerLogResponse, error) {
	return nil, fmt.Errorf(" does not support reopening container log files")
}

// GetContainerLogs get container logs directly from docker daemon.
func (ds *templateService) GetContainerLogs(
	_ context.Context,
	pod *v1.Pod,
	containerID config.ContainerID,
	logOptions *v1.PodLogOptions,
	stdout, stderr io.Writer,
) error {

	return nil
}

// GetContainerLogTail attempts to read up to MaxContainerTerminationMessageLogLength
// from the end of the log when docker is configured with a log driver other than json-log.
// It reads up to MaxContainerTerminationMessageLogLines lines.
func (ds *templateService) GetContainerLogTail(
	uid config.UID,
	name, namespace string,
	containerID config.ContainerID,
) (string, error) {

	return "mock log", nil
}

// criSupportedLogDrivers are log drivers supported by native CRI integration.
var criSupportedLogDrivers = []string{"json-file"}

// IsCRISupportedLogDriver checks whether the logging driver used by docker is
// supported by native CRI integration.
func (ds *templateService) IsCRISupportedLogDriver() (bool, error) {

	return false, nil
}
