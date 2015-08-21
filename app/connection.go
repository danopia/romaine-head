package app

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/danopia/romaine-head/head"
)

func HandleConn(conn *websocket.Conn) {
	log.Println("Client app connected")
	head.CurrentApp = conn

	for {
		var request Request
		if err := conn.ReadJSON(&request); err != nil {
			log.Println("Error reading app json: ", err)
			return
		}

		response := HandleRequest(&request)
		conn.WriteJSON(response)
	}
}
