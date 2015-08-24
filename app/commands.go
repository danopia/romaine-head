package app

import (
	"log"
	"strings"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/head"
)

func listRoots() []map[string]interface{} {
	output, _ := common.RunCmd("ls", "-1", "/mnt/stateful_partition/crouton/chroots")
	names := strings.Split(output, "\n")
	chroots := make([]map[string]interface{}, len(names))

	for i, name := range names {
		chroot := make(map[string]interface{})
		chroot["key"] = name

		if val, ok := head.GetLeaf(name); ok {
			chroot["state"] = val.State
		} else {
			chroot["state"] = "stopped"
		}

		chroots[i] = chroot
	}

	return chroots
}

func runCrouton(args []string) string {
	output, _ := common.RunCmd("sh", append([]string{"/home/chronos/user/Downloads/crouton"}, args...)...)
	return output
}

func buildChroot(args []string, stdin string) string {
	output, _ := common.RunCmdWithStdin("sh", append([]string{"/home/chronos/user/Downloads/crouton"}, args...), stdin)
	return output
}

func runInChroot(chroot string, cmd []string, context string) {
	// output, _ := common.RunCmd("sudo", append([]string{"enter-chroot"}, cmd...)...)

	if leaf, ok := head.GetLeaf(chroot); ok {
		if leaf.Conn != nil {
			leaf.Conn.WriteJSON(&common.Packet{
				Cmd:     "exec",
				Context: context,
				Extras: map[string]interface{}{
					"Path": cmd[0],
					"Args": cmd[1:],
				},
			})

		} else {
			log.Printf("chroot %s isn't running", chroot)
		}
	}
}
