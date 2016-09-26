package common

import (
	"log"
	"os/exec"
	"strings"
	"syscall"
)

func RunCmd(name string, arg ...string) (string, int) {
	cmd := exec.Command(name, arg...)

	output, err := cmd.CombinedOutput()
	if _, ok := err.(*exec.ExitError); !ok && err != nil {
		log.Printf("cmd error %v\n", err)
	}

	trimmedOutput := strings.TrimRight(string(output), "\n")
	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	return trimmedOutput, waitStatus.ExitStatus()
}

func RunCmdWithStdin(name string, arg []string, stdin string) (string, int) {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = strings.NewReader(stdin)

	output, err := cmd.CombinedOutput()
	if _, ok := err.(*exec.ExitError); !ok && err != nil {
		log.Printf("cmd error %v\n", err)
	}

	trimmedOutput := strings.TrimRight(string(output), "\n")
	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	return trimmedOutput, waitStatus.ExitStatus()
}
