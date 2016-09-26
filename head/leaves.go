package head

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/ddp"
	"github.com/kr/pty"
	"github.com/kr/text"
)

const headPath = "~/Downloads/romaine-head.run"

var Leaves = make(map[string]*Leaf)

func GetLeaf(leaf string) (val *Leaf, ok bool) {
	val, ok = Leaves[leaf]
	return
}

// master process listens on this port
// slave chroot processes connect to the master for instructions
var port = 6206

func StartLeaf(leaf string, password string) *Leaf {
	if entry, ok := GetLeaf(leaf); ok && entry.Anchor != nil {
		return entry
	}

	secret := common.GenerateSecret()

	log.Printf("Starting %s under port %d", leaf, port)
	command := fmt.Sprintf("%s -- --mode leaf --port %d --secret %s 2>&1", headPath, port, secret)

	prefix := []byte(fmt.Sprintf("[%s] ", leaf))
	output := text.NewIndentWriter(os.Stdout, prefix)

	entry := &Leaf{
		Id:     leaf,
		State:  "launching",
		Secret: secret,
		Source: make(chan common.Packet),
		Sink:   make(chan common.Packet),
		Anchor: exec.Command("enter-chroot", "-n", leaf, "sh", "-c", command),
	}
	Leaves[leaf] = entry
	ddp.Chroots.SetField(leaf, "status", entry.State)

	f, err := pty.Start(entry.Anchor)
	if err != nil {
		// TODO: enough cleanup?
		log.Printf("Leaf %s PTY failed with %+v", leaf, err)
		entry.State = "crashed"
		entry.Anchor = nil
		ddp.Chroots.SetField(leaf, "status", entry.State)

	} else {
		entry.Pty = f
		go io.Copy(output, entry.Pty)

		if password != "" {
			// TODO: don't write this until prompted
			f.Write([]byte(password + "\n"))
		}
	}

	go func() {
		err := entry.Anchor.Wait()
		log.Printf("Leaf %s exited with %+v", leaf, err)

		if err != nil {
			entry.State = "crashed"
		} else {
			entry.State = "stopped"
		}
		entry.Anchor = nil
		entry.Pty = nil

		ddp.Chroots.SetField(leaf, "status", entry.State)
	}()

	return entry
}
