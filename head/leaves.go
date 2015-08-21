package head

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/danopia/romaine-head/common"
)

var leaves = make(map[string]*Leaf)

func GetLeaf(leaf string) (val *Leaf, ok bool) {
	val, ok = leaves[leaf]
	return
}

var port = 6205

func StartLeaf(leaf string) *Leaf {
	if entry, ok := GetLeaf(leaf); ok {
		return entry
	}

	secret := common.GenerateSecret()

	log.Printf("Starting %s under port %d", leaf, port)
	command := fmt.Sprintf("~/Downloads/romaine-head.run -- --mode leaf --port %d --secret %s 2>&1 | nc localhost 5000", port, secret)

	entry := &Leaf{
		State:  "launching",
		Secret: secret,
		Anchor: exec.Command("enter-chroot", "fish", "-c", command),
	}
	entry.Anchor.Start()

	go func() {
		err := entry.Anchor.Wait()
		log.Printf("Leaf %s exited with %+v", leaf, err)

		if err != nil {
			entry.State = "crashed"
		} else {
			entry.State = "stopped"
		}
	}()

	leaves[leaf] = entry
	return entry
}
