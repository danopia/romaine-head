package stem

import (
	"log"

	"github.com/danopia/romaine-head/common"
	// "github.com/danopia/romaine-head/head"
	"github.com/gorilla/websocket"
)

// HandleLeafConn and run the Head program
// Starting with auth
func HandleLeafConn(conn *websocket.Conn) {
	log.Println("Leaf connected")

	leaf, ok := AuthLeafConn(conn)
	if !ok {
		conn.Close()
		return
	}

	go leaf.PumpSink()
	leaf.Sink <- common.Packet{
		Cmd: "ready",
	}

	for {
		var packet common.Packet
		if err := conn.ReadJSON(&packet); err != nil {
			log.Println("Error reading leaf json: ", err)
			return
		}

		HandlePacket(&packet, leaf)
	}
}
