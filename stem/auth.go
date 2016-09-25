package stem

import (
	"log"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/ddp"
	"github.com/danopia/romaine-head/head"
	"github.com/gorilla/websocket"
)

// AuthLeafConn using a secret token payload
// TODO: all this error handling lol
func AuthLeafConn(conn *websocket.Conn) (*head.Leaf, bool) {
	var packet common.Packet
	if err := conn.ReadJSON(&packet); err != nil {
		log.Println("Error reading leaf json: ", err)
		return nil, false
	}
	log.Printf("leaf auth <<< %+v\n", packet)

	if packet.Cmd != "auth" {
		return nil, false
	}

	leaf, ok := leafBySecret(packet.Context)
	if !ok {
		return nil, false
	}

	log.Printf("Leaf identified as %s", leaf.Id)

	leaf.Conn = conn
	leaf.State = "running"
	ddp.Chroots.SetField(leaf.Id, "status", "running")

	return leaf, true
}

func leafBySecret(secret string) (*head.Leaf, bool) {
	for _, leaf := range head.Leaves {
		if leaf.Secret == secret {
			return leaf, true
		}
	}

	return nil, false
}
