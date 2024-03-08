package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

// the next line is just for reference and contains the os.Args slice
// [/tmp/tmp.pPHgBn run codecraftersio/docker-challenge /usr/local/bin/docker-explorer echo 13254]

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	EnterNewJail(os.Args[3])

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	cmd.Run()
	os.Exit(cmd.ProcessState.ExitCode())
}

func EnterNewJail(filepath string) {
	tempDirPath, err := os.MkdirTemp("", "temp_folder_*")
	if err != nil {
		throwError(err, "Unable to create temp directory")
	}

	defer os.Remove(tempDirPath)

	err = os.Chmod(tempDirPath, 0777)
	if err != nil {
		throwError(err, "Error modifying rwx on tempDirPath")
	}

	filepathSplit := strings.Split(filepath, "/")
	jailedDirPath := path.Join(tempDirPath, strings.Join(filepathSplit[0:len(filepathSplit)-1], "/"))

	err = os.MkdirAll(jailedDirPath, 0777)
	if err != nil {
		throwError(err, fmt.Sprintf("Error in mkdirall to %s", jailedDirPath))
	}

	// not in the mood to write code to copy files
	err = os.Link(filepath, path.Join(tempDirPath, filepath))
	if err != nil {
		throwError(err, "Error in link command")
	}

	err = syscall.Chroot(tempDirPath)
	if err != nil {
		throwError(err, fmt.Sprintf("Error in executing chroot on %s", tempDirPath))
	}
}

func throwError(err error, msg string) {
	log.Fatalf("%s: %v\n", msg, err)
	os.Exit(1)
}
