package leaf

import (
	"fmt"
	"log"

	"github.com/danopia/romaine-head/common"
	"github.com/gorilla/websocket"
)

func ConnectToHead(url string, secret string) {
	var cstDialer = websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, _, err := cstDialer.Dial(url, nil)
	if err != nil {
		fmt.Printf("Dial: %v", err)
	}

	log.Printf("ok")

	// Send a welcome packet for auth
	conn.WriteJSON(&common.Packet{
		Cmd:     "auth",
		Context: secret,
	})

	for {
		var p common.Packet

		if err := conn.ReadJSON(&p); err != nil {
			log.Println("Error reading head's json: ", err)
			return
		}

		response := HandlePacket(p)

		if response != nil {
			conn.WriteJSON(response)
		}
	}
}
