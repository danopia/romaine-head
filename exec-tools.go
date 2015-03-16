package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func runCmd(name string, arg ...string) (string, int) {
	cmd := exec.Command(name, arg...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd error %s\n", err)
		//panic(err)
		// TODO return the exit code from err
	}

	return strings.TrimRight(string(output), "\n"), -1
}

func listRoots() string {
	output, _ := runCmd("ls", "-1", "/mnt/stateful_partition/crouton/chroots")
	return output
}

func runCrouton(args []string) string {
	output, _ := runCmd("sh", append([]string{"/home/chronos/user/Downloads/crouton"}, args...)...)
	return output
}

func runInChroot(chroot string, cmd []string) string {
	output, _ := runCmd("sudo", append([]string{"enter-chroot"}, cmd...)...)
	return output
}
