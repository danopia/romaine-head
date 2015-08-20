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
	var welcome = common.Packet{
		Cmd:     "auth",
		Context: secret,
	}
	conn.WriteJSON(&welcome)

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

	/*
		url := "ws://localhost:6205/ws"
		ws, err := websocket.Dial(url, "", "http://localhost")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
			log.Fatal(err)
		}
		var msg = make([]byte, 512)
		var n int
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received: %s.\n", msg[:n])
	*/
}