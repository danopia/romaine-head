package stem

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/danopia/romaine-head/common"
)

func HandleLeafConn(conn *websocket.Conn) {
	log.Println("Leaf connected")

	for {
		var packet common.Packet
		if err := conn.ReadJSON(&packet); err != nil {
			log.Println("Error reading leaf json: ", err)
			return
		}

		HandlePacket(&packet, conn)
	}
}
