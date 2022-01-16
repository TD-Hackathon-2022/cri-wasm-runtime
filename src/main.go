package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/diannaowa/cri-template/cmd"

	"k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := cmd.NewTemplateCRICommand(server.SetupSignalHandler())
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
