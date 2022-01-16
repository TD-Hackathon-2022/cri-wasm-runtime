package options

import (
	"github.com/diannaowa/cri-template/config"
	"github.com/spf13/pflag"
	"runtime"
)

// TemplateCRIFlags contains configuration flags for cri-template
type TemplateCRIFlags struct {
	config.ContainerRuntimeOptions
	// remoteRuntimeEndpoint is the endpoint of backend runtime service
	RemoteRuntimeEndpoint string
}

// NewTemplateCRIFlags will create a new TemplateCRIFlags with default values
func NewTemplateCRIFlags() *TemplateCRIFlags {
	remoteRuntimeEndpoint := ""
	if runtime.GOOS == "linux" {
		remoteRuntimeEndpoint = "unix:///var/run/cri-template.sock"
	} else if runtime.GOOS == "windows" {
		remoteRuntimeEndpoint = "npipe:////./pipe/cri-template"
	}

	return &TemplateCRIFlags{
		RemoteRuntimeEndpoint: remoteRuntimeEndpoint,
	}
}

// TemplateCRIServer encapsulates all of the parameters necessary for starting up
// a kubelet. These can either be set via command line or directly.
type TemplateCRIServer struct {
	TemplateCRIFlags
}

// AddFlags adds flags for a specific TemplateCRIServer to the specified FlagSet
func (s *TemplateCRIServer) AddFlags(fs *pflag.FlagSet) {
	s.TemplateCRIFlags.AddFlags(fs)
}

// AddFlags adds flags for a specific TemplateCRIFlags to the specified FlagSet
func (f *TemplateCRIFlags) AddFlags(mainfs *pflag.FlagSet) {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	defer func() {
		// Un-hide deprecated flags. We want deprecated flags to show in cri-dockerd's help.
		// We have some hidden flags, but we may as well un-hide these when they are deprecated,
		// as silently deprecating and removing (even hidden) things is unkind to people who use them.
		fs.VisitAll(func(f *pflag.Flag) {
			if len(f.Deprecated) > 0 {
				f.Hidden = false
			}
		})
		mainfs.AddFlagSet(fs)
	}()

	fs.StringVar(
		&f.RemoteRuntimeEndpoint,
		"container-runtime-endpoint",
		f.RemoteRuntimeEndpoint,
		"The endpoint of backend runtime service. Currently unix socket and tcp endpoints are supported on Linux, while npipe and tcp endpoints are supported on windows.  Examples:'unix:///var/run/cri-template.sock', 'npipe:////./pipe/cri-template'",
	)
}
