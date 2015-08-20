package head

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/danopia/romaine-head/common"
)

func HandleLeafStem(p common.Packet, conn *websocket.Conn) {
	log.Printf("leaf <<< %+v\n", p)

	switch p.Cmd {
	case "list chroots":

	case "start chroot":

	case "run crouton":

	case "run in chroot":

	default:
		log.Fatal("Leaf sent unknown packet " + p.Cmd)
	}

	//log.Printf(">>> response to %s: %+v\n", p.Context, response)
	return
}
