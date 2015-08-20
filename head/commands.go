package head

import (
	"strings"

	"github.com/danopia/romaine-head/common"
)

func listRoots() []map[string]interface{} {
	output, _ := common.RunCmd("ls", "-1", "/mnt/stateful_partition/crouton/chroots")
	names := strings.Split(output, "\n")
  chroots := make([]map[string]interface{}, len(names))

  for i, name := range names {
    chroot := make(map[string]interface{})
    chroot["key"] = name

    if val, ok := getLeaf(name); ok {
			chroot["running"] = val.Running
      //chroot["port"] = val.Port
    }

    chroots[i] = chroot
  }

	return chroots
}

func runCrouton(args []string) string {
	output, _ := common.RunCmd("sh", append([]string{"/home/chronos/user/Downloads/crouton"}, args...)...)
	return output
}

func runInChroot(chroot string, cmd []string) string {
	output, _ := common.RunCmd("sudo", append([]string{"enter-chroot"}, cmd...)...)
	return output
}
