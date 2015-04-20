package head

import (
	"log"
	"fmt"
  "os/exec"
)

var leaves = make(map[string]*Leaf)

func getLeaf(leaf string) (val *Leaf, ok bool) {
  val, ok = leaves[leaf]
  return
}

var nextPort = 6206
func startLeaf(leaf string) *Leaf {
	if entry, ok := getLeaf(leaf); ok {
		return entry
	}

  log.Printf("Starting %s on port %d", leaf, nextPort)
	command := fmt.Sprintf("romaine-head --mode leaf --port %d", nextPort)

	entry := &Leaf{
    Running: true,
    Port:    nextPort,
    Anchor:  exec.Command("enter-chroot", "fish", "-c", command),
  }
	entry.Anchor.Start()

  nextPort += 1
  leaves[leaf] = entry
	return entry
}
