package core

import (
	"context"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"k8s.io/client-go/tools/remotecommand"
)

// NativeExecHandler executes commands in Docker containers using Docker's exec API.
type NativeExecHandler struct {
	templateService *templateService
}

// ExecInContainer executes the cmd in container using the Docker's exec API
func (ne *NativeExecHandler) ExecInContainer(
	ctx context.Context,
	containerID string,
	cmd []string,
	stdin io.Reader,
	stdout, stderr io.WriteCloser,
	tty bool,
	resize <-chan remotecommand.TerminalSize,
	timeout time.Duration,
) error {
	baseWasmDir := "/root/wasm/hello"
	containerCache := ne.templateService.containerCache[containerID]
	if containerCache == nil {
		return fmt.Errorf("cannot find container")
	}

	// imageName -> wasm file path
	imageName := ne.templateService.imageCache[containerCache.config.GetImage().GetImage()].name
	index := strings.Index(imageName, ":")
	var wasmName = ""
	if index != -1 {
		wasmName = imageName[:index] + ".wasm"
	} else {
		wasmName = imageName + ".wasm"
	}

	wasmOpenDir := "/var/run/wasm/" + containerCache.id + "/" + uuid.New()
	err := os.MkdirAll(wasmOpenDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("exec error, create wasmOpenDir failed")
	}

	wasmFilePath := baseWasmDir + "/" + wasmName

	args := strings.Join(cmd, " ")
	outputFilePath := wasmOpenDir + "/result"
	_, createErr := os.Create(outputFilePath)
	if createErr != nil {
		return fmt.Errorf("exec error, create result file failed")
	}
	// imageName -> wasm file path

	commandStr := fmt.Sprintf("wasmtime run --dir %s %s %s %s", wasmOpenDir, wasmFilePath, args, outputFilePath)
	command := exec.Command("/bin/bash", "-c", commandStr)
	output, execErr := command.CombinedOutput()
	if execErr != nil {
		fmt.Println("run error, error: ", string(output))
		stderr.Write(output)
	} else {
		resultBytes, readResultErr := ioutil.ReadFile(outputFilePath)
		if readResultErr != nil {
			return fmt.Errorf("exec error,read result file failed")
		}
		stdout.Write(resultBytes)
	}

	removeFileErr := os.RemoveAll(wasmOpenDir)
	if removeFileErr != nil {
		logrus.Infof("exec error,remove wasmOpenDir failed: %s", wasmOpenDir)
	}

	logrus.Infof("exec container wasmName: %s", wasmName)
	return nil
}
