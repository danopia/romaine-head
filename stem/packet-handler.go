package stem

import (
	"log"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/head"
	"github.com/gorilla/websocket"
)

func HandlePacket(p *common.Packet, conn *websocket.Conn) {
	log.Printf("leaf <<< %+v\n", p)

	switch p.Cmd {

	// Authenticate payload with a secret token
	case "auth":
		for name, leaf := range head.Leaves {
			if leaf.Secret == p.Context {
				leaf.Conn = conn
				leaf.State = "running"

				log.Printf("Leaf identified as %s", name)
				return
			}
		}

	// Response from an execution
	case "exec":
		head.CurrentApp.WriteJSON(&map[string]interface{}{
			"context": p.Context,
			"output":  p.Extras["Output"].(string),
		})

	default:
		log.Printf("Leaf sent unknown packet {}", p.Cmd)
	}

	//log.Printf(">>> response to %s: %+v\n", p.Context, response)
}
