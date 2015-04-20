package common

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunCmd(name string, arg ...string) (string, int) {
	cmd := exec.Command(name, arg...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd error %s\n", err)
		//panic(err)
		// TODO return the exit code from err
	}

	return strings.TrimRight(string(output), "\n"), -1
}
