package app

import (
	"os"
	"log"
	"strings"
	"bufio"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/head"
)

const rootPath = "/mnt/stateful_partition/crouton/chroots"

func listRoots() []map[string]interface{} {
	output, _ := common.RunCmd("ls", "-1", rootPath)
	names := strings.Split(output, "\n")
	chroots := make([]map[string]interface{}, len(names))

	for i, name := range names {
		chrootPath := rootPath + "/" + name

		chroot := make(map[string]interface{})
		chroot["key"] = name

		// is there an anchor running already?
		if val, ok := head.GetLeaf(name); ok {
			chroot["state"] = val.State
		} else {
			chroot["state"] = "stopped"
		}

		// which targets does the chroot include?
	  lines, err := readLines(chrootPath + "/.crouton-targets")
	  if err != nil {
	    log.Println("[chronos] failed to read targets from %s : %s", name, err)
	  } else {
			chroot["targets"] = lines
		}

		// is the chroot encrypted?
		_, ecryptErr := os.Stat(chrootPath + "/.ecryptfs")
		chroot["encrypted"] = !os.IsNotExist(ecryptErr)

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

// readLines reads a whole file into memory
// and returns a slice of its lines.
// http://stackoverflow.com/a/18479916
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}
