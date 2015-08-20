package head

import (
	"log"

	"github.com/danopia/romaine-head/common"
	"github.com/gorilla/websocket"
)

func HandleLeafStem(p common.Packet, conn *websocket.Conn) {
	log.Printf("leaf <<< %+v\n", p)

	switch p.Cmd {

	// Authenticate payload with a secret token
	case "auth":
		for name, leaf := range leaves {
			if leaf.Secret == p.Context {
				leaf.Conn = conn
				leaf.State = "running"

				log.Printf("Leaf identified as %s", name)
				return
			}
		}

	// Response from an execution
	case "exec":
		common.CurrentApp.WriteJSON(&map[string]interface{}{
			"context": p.Context,
			"output":  p.Extras["Output"].(string),
		})

	default:
		log.Printf("Leaf sent unknown packet {}", p.Cmd)
	}

	//log.Printf(">>> response to %s: %+v\n", p.Context, response)
}
