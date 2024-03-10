//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/codecrafters-io/docker-starter-go/dockerutils"
	"github.com/codecrafters-io/docker-starter-go/throwerror"
)

// the next line is just for reference and contains the os.Args slice
// [/tmp/tmp.pPHgBn run codecraftersio/docker-challenge /usr/local/bin/docker-explorer echo 13254]

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	image := os.Args[2]

	jailPath := EnterNewJail(os.Args[3], image)
	defer os.Remove(jailPath)

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID,
	}

	cmd.Run()
	os.Exit(cmd.ProcessState.ExitCode())
}

func EnterNewJail(filepath string, image string) string {
	tempDirPath, err := os.MkdirTemp("", "temp_folder_*")
	if err != nil {
		throwerror.ThrowError(err, "Unable to create temp directory")
	}

	err = os.Chmod(tempDirPath, 0777)
	if err != nil {
		throwerror.ThrowError(err, "Error modifying rwx on tempDirPath")
	}

	authToken := dockerutils.GetAuthToken(image)
	manifest := dockerutils.GetManifest(image, authToken)
	dockerutils.DownloadAndExtractLayers(manifest.Layers, image, authToken, tempDirPath)

	err = syscall.Chroot(tempDirPath)
	if err != nil {
		throwerror.ThrowError(err, fmt.Sprintf("Error in executing chroot on %s", tempDirPath))
	}

	return tempDirPath
}
