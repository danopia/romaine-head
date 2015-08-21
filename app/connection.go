package app

import (
	"log"

	"github.com/danopia/romaine-head/head"
	"github.com/gorilla/websocket"
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
