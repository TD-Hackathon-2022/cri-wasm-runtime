package cmd

import (
	"fmt"
	"github.com/diannaowa/cri-template/backend"
	"github.com/diannaowa/cri-template/cmd/cri/options"
	"github.com/diannaowa/cri-template/core"
	"github.com/diannaowa/cri-template/streaming"
	"github.com/diannaowa/cri-template/version"
	"github.com/sirupsen/logrus"

	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	componentTemplateCRI = "cri-template"
)

// NewTemplateCRICommand creates a *cobra.Command object with default parameters
func NewTemplateCRICommand(stopCh <-chan struct{}) *cobra.Command {
	cleanFlagSet := pflag.NewFlagSet(componentTemplateCRI, pflag.ContinueOnError)
	kubeletFlags := options.NewTemplateCRIFlags()

	cmd := &cobra.Command{
		Use:  componentTemplateCRI,
		Long: `CRI that connects to the Docker Daemon`,
		// cri-dockerd has special flag parsing requirements to enforce flag precedence rules,
		// so we do all our parsing manually in Run, below.
		// DisableFlagParsing=true provides the full set of flags passed to cri-dockerd in the
		// `args` arg to Run, without Cobra's interference.
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			// initial flag parse, since we disable cobra's flag parsing
			if err := cleanFlagSet.Parse(args); err != nil {
				cmd.Usage()
				logrus.Fatal(err)
			}

			// check if there are non-flag arguments in the command line
			cmds := cleanFlagSet.Args()
			if len(cmds) > 0 {
				cmd.Usage()
				logrus.Fatalf("Unknown command: %s", cmds[0])
			}

			// short-circuit on help
			help, err := cleanFlagSet.GetBool("help")
			if err != nil {
				logrus.Fatal(`"help" flag is non-bool`)
			}
			if help {
				cmd.Help()
				return
			}

			verflag, _ := cleanFlagSet.GetBool("version")
			if verflag {
				fmt.Fprintf(
					cmd.OutOrStderr(),
					"%s %s\n",
					version.PlatformName,
					version.FullVersion(),
				)
				return
			}

			infoflag, _ := cleanFlagSet.GetBool("buildinfo")
			if infoflag {
				fmt.Fprintf(
					cmd.OutOrStderr(),
					"Program: %s\nVersion: %s\nBuildTime: %s\nGitCommit: %s\n",
					version.PlatformName,
					version.FullVersion(),
					version.BuildTime,
					version.GitCommit,
				)
				return
			}

			logFlag, _ := cleanFlagSet.GetString("log-level")
			if logFlag != "" {
				level, err := logrus.ParseLevel(logFlag)
				if err != nil {
					logrus.Fatalf("Unknown log level: %s", logFlag)
				}
				logrus.SetLevel(level)
			}

			if err := RunCriTemplate(kubeletFlags, stopCh); err != nil {
				logrus.Fatal(err)
			}
		},
	}

	// keep cleanFlagSet separate, so Cobra doesn't pollute it with the global flags
	kubeletFlags.AddFlags(cleanFlagSet)
	cleanFlagSet.BoolP("help", "h", false, fmt.Sprintf("Help for %s", cmd.Name()))
	cleanFlagSet.Bool("version", false, "Prints the version of cri-dockerd")
	cleanFlagSet.Bool("buildinfo", false, "Prints the build information about cri-dockerd")
	cleanFlagSet.String("log-level", "info", "The log level for cri-docker")

	// ugly, but necessary, because Cobra's default UsageFunc and HelpFunc pollute the flagset with global flags
	const usageFmt = "Usage:\n  %s\n\nFlags:\n%s"
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine(), cleanFlagSet.FlagUsagesWrapped(2))
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(
			cmd.OutOrStdout(),
			"%s\n\n"+usageFmt,
			cmd.Long,
			cmd.UseLine(),
			cleanFlagSet.FlagUsagesWrapped(2),
		)
	})
	return cmd
}

// RunCriTemplate starts cri-dockerd
func RunCriTemplate(f *options.TemplateCRIFlags, stopCh <-chan struct{}) error {
	r := &f.ContainerRuntimeOptions

	// Initialize streaming configuration. (Not using TLS now)
	streamingConfig := &streaming.Config{
		// Use a relative redirect (no scheme or host).
		BaseURL:                         &url.URL{Path: "/cri/"},
		StreamIdleTimeout:               r.StreamingConnectionIdleTimeout.Duration,
		StreamCreationTimeout:           streaming.DefaultConfig.StreamCreationTimeout,
		SupportedRemoteCommandProtocols: streaming.DefaultConfig.SupportedRemoteCommandProtocols,
		SupportedPortForwardProtocols:   streaming.DefaultConfig.SupportedPortForwardProtocols,
	}

	// Standalone cri-template will always start the local streaming backend.
	ds, err := core.NewTemplateService(
		streamingConfig,
		r.CriTemplateRootDirectory,
	)
	if err != nil {
		return err
	}

	logrus.Info("Starting the GRPC backend for the Template CRI interface.")
	server := backend.NewCriTemplateServer(f.RemoteRuntimeEndpoint, ds)
	if err := server.Start(); err != nil {
		return err
	}

	<-stopCh
	return nil
}
