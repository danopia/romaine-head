package head

import (
	"os/exec"

	"github.com/danopia/romaine-head/common"
	"github.com/gorilla/websocket"
)

// Leaf corresponding to a single chroot
// And hopefully a process anchor
type Leaf struct {
	Id     string
	State  string
	Secret string
	Anchor *exec.Cmd
	Conn   *websocket.Conn
  Source chan common.Packet
  Sink   chan common.Packet
}

// PumpSink into the websocket
func (l *Leaf) PumpSink() {
	for packet := range l.Sink {
		l.Conn.WriteJSON(packet)
		// TODO: write error handling,
		// close the channels
	}
}
