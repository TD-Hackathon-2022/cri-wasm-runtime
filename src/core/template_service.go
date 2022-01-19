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

package core

import (
	"context"
	"github.com/diannaowa/cri-template/config"
	"github.com/diannaowa/cri-template/metrics"
	"github.com/diannaowa/cri-template/streaming"
	"github.com/sirupsen/logrus"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/clock"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	v1cri "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"net/http"
	"os"
)

const (
	templateRuntimeName = "cri-template"
	kubeAPIVersion      = "0.1.0"

	// String used to detect docker host mode for various namespaces (e.g.
	// networking). Must match the value returned by docker inspect -f
	// '{{.HostConfig.NetworkMode}}'.
	namespaceModeHost = "host"
	maxMsgSize        = 1024 * 1024 * 16
)

// CRIService includes all methods necessary for a CRI backend.
type CRIService interface {
	runtimeapi.RuntimeServiceServer
	runtimeapi.ImageServiceServer
	Start() error
}

// DockerService is an interface that embeds the new RuntimeService and
// ImageService interfaces.
type DockerService interface {
	CRIService

	// For serving streaming calls.
	http.Handler

	// GetContainerLogs gets logs for a specific container.
	GetContainerLogs(
		context.Context,
		*v1.Pod,
		config.ContainerID,
		*v1.PodLogOptions,
		io.Writer,
		io.Writer,
	) error

	// IsCRISupportedLogDriver checks whether the logging driver used by docker is
	// supported by native CRI integration.
	IsCRISupportedLogDriver() (bool, error)

	// Get the last few lines of the logs for a specific container.
	GetContainerLogTail(
		uid config.UID,
		name, namespace string,
		containerID config.ContainerID,
	) (string, error)
}

// NewTemplateService creates a new `DockerService` struct.
func NewTemplateService(
	streamingConfig *streaming.Config,
	criTemplateRootDir string,
) (DockerService, error) {

	ds := &templateService{
		os: config.RealOS{},

		clock:          clock.RealClock{},
		sandboxCache:   make(map[string]*sandboxCacheModel),
		imageCache:     make(map[string]*imageCacheModel),
		containerCache: make(map[string]*containerCacheModel),
	}

	ds.streamingRuntime = &streaming.StreamingRuntime{
		ExecHandler: &NativeExecHandler{
			templateService: ds,
		},
	}

	// create streaming backend if configured.
	if streamingConfig != nil {
		var err error
		ds.streamingServer, err = streaming.NewServer(*streamingConfig, ds.streamingRuntime)
		if err != nil {
			return nil, err
		}
	}

	// Register prometheus metrics.
	metrics.Register()
	return ds, nil
}

type sandboxCacheModel struct {
	id             string
	config         *v1cri.PodSandboxConfig
	status         *v1cri.PodSandboxStatus
	containerIdMap map[string]string
}

type imageCacheModel struct {
	id            string
	name          string
	image         *v1cri.ImageSpec
	imageStatus   *v1cri.Image
	sandboxConfig *v1cri.PodSandboxConfig
}

type containerCacheModel struct {
	id            string
	config        *v1cri.ContainerConfig
	status        *v1cri.ContainerStatus
	sandboxId     string
	sandboxConfig *v1cri.PodSandboxConfig
}

type templateService struct {
	os               config.OSInterface
	streamingRuntime *streaming.StreamingRuntime
	streamingServer  streaming.Server

	clock          clock.RealClock
	sandboxCache   map[string]*sandboxCacheModel
	imageCache     map[string]*imageCacheModel
	containerCache map[string]*containerCacheModel
}

// Version returns the runtime name, runtime version and runtime API version
func (ds *templateService) Version(
	_ context.Context,
	r *runtimeapi.VersionRequest,
) (*runtimeapi.VersionResponse, error) {
	return &runtimeapi.VersionResponse{
		Version:           kubeAPIVersion,
		RuntimeName:       templateRuntimeName,
		RuntimeVersion:    "0.1.0",
		RuntimeApiVersion: "0.0",
	}, nil
}

// UpdateRuntimeConfig updates the runtime config. Currently only handles podCIDR updates.
func (ds *templateService) UpdateRuntimeConfig(
	_ context.Context,
	r *runtimeapi.UpdateRuntimeConfigRequest,
) (*runtimeapi.UpdateRuntimeConfigResponse, error) {

	return &runtimeapi.UpdateRuntimeConfigResponse{}, nil
}

// Start initializes and starts components in templateService.
func (ds *templateService) Start() error {
	ds.initCleanup()
	go func() {
		if err := ds.streamingServer.Start(true); err != nil {
			logrus.Errorf("Streaming backend stopped unexpectedly: %v", err)
			os.Exit(1)
		}
	}()

	return nil
}

// Status returns the status of the runtime.
func (ds *templateService) Status(
	_ context.Context,
	r *runtimeapi.StatusRequest,
) (*runtimeapi.StatusResponse, error) {
	runtimeReady := &runtimeapi.RuntimeCondition{
		Type:   runtimeapi.RuntimeReady,
		Status: true,
	}
	networkReady := &runtimeapi.RuntimeCondition{
		Type:   runtimeapi.NetworkReady,
		Status: true,
	}
	conditions := []*runtimeapi.RuntimeCondition{runtimeReady, networkReady}
	// if something not ready
	//
	status := &runtimeapi.RuntimeStatus{Conditions: conditions}
	return &runtimeapi.StatusResponse{Status: status}, nil
}

func (ds *templateService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ds.streamingServer != nil {
		ds.streamingServer.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// initCleanup is responsible for cleaning up any crufts left by previous
// runs. If there are any errors, it simply logs them.
func (ds *templateService) initCleanup() {
}
