package leaf

import (
	"fmt"
	"log"

	"github.com/danopia/romaine-head/common"
	"github.com/gorilla/websocket"
)

type Stalk struct {
  Conn *websocket.Conn
  Source chan common.Packet
  Sink chan common.Packet
}

func (s *Stalk) pumpSink() {
	for message := range s.Sink {
		s.Conn.WriteJSON(message)
	}
}

func ConnectToHead(url string, secret string) {
	var cstDialer = websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, _, err := cstDialer.Dial(url, nil)
	if err != nil {
		fmt.Printf("Dial: %v", err)
	}

	s := Stalk{
		Conn: conn,
		Source: make(chan common.Packet),
		Sink: make(chan common.Packet),
	}

	log.Printf("Connection established to romaine-head")
	go s.pumpSink()

	// Send a welcome packet for auth
	s.Sink <- common.Packet{
		Cmd:     "auth",
		Context: secret,
	}

	for { // each incoming
		var p common.Packet

		if err := conn.ReadJSON(&p); err != nil {
			log.Println("Error reading head's json: ", err)
			return
		}

		s.handlePacket(p)

		// TODO: handle this better
		if p.Cmd == "shutdown" {
			break
		}
	}
}
