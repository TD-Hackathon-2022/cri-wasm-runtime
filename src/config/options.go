package config

import (
	"github.com/spf13/pflag"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ContainerRuntimeOptions contains runtime options
type ContainerRuntimeOptions struct {
	// General options.

	// CriTemplateRootDirectory is the path to the cri-template root directory. Defaults to
	// /var/lib/cri-template if unset. Exposed for integration testing (e.g. in OpenShift).
	CriTemplateRootDirectory string
	RuntimeRequestTimeout    v1.Duration
	// streamingConnectionIdleTimeout is the maximum time a streaming connection
	// can be idle before the connection is automatically closed.
	StreamingConnectionIdleTimeout v1.Duration
}

// AddFlags has the set of flags needed by cri-template
func (s *ContainerRuntimeOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&s.CriTemplateRootDirectory,
		"cri-dockerd-root-directory",
		s.CriTemplateRootDirectory,
		"Path to the cri-dockerd root directory.",
	)

}
