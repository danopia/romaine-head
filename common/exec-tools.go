package common

import (
	"log"
	"os/exec"
	"strings"
)

func RunCmd(name string, arg ...string) (string, int) {
	cmd := exec.Command(name, arg...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd error %s\n", err)
		//panic(err)
		// TODO return the exit code from err
	}

	return strings.TrimRight(string(output), "\n"), -1
}

func RunCmdWithStdin(name string, arg []string, stdin string) (string, int) {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = strings.NewReader(stdin)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd error %s\n", err)
	}

	return strings.TrimRight(string(output), "\n"), -1
}
