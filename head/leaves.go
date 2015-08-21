package head

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/danopia/romaine-head/common"
	"github.com/kr/text"
)

const headPath = "~/Downloads/romaine-head.run"

var Leaves = make(map[string]*Leaf)

func GetLeaf(leaf string) (val *Leaf, ok bool) {
	val, ok = Leaves[leaf]
	return
}

var port = 6205

func StartLeaf(leaf string) *Leaf {
	if entry, ok := GetLeaf(leaf); ok {
		return entry
	}

	secret := common.GenerateSecret()

	log.Printf("Starting %s under port %d", leaf, port)
	command := fmt.Sprintf("%s -- --mode leaf --port %d --secret %s 2>&1", headPath, port, secret)

	prefix := []byte(fmt.Sprintf("[%s] ", leaf))
	output := text.NewIndentWriter(os.Stdout, prefix)

	entry := &Leaf{
		State:  "launching",
		Secret: secret,
		Anchor: exec.Command("enter-chroot", "sh", "-c", command),
	}
	entry.Anchor.Stdout = output
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

	Leaves[leaf] = entry
	return entry
}
