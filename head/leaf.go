package head

import (
	"os/exec"

	"github.com/gorilla/websocket"
)

type Leaf struct {
	State  string
	Secret string
	Anchor *exec.Cmd
	Conn   *websocket.Conn
}
