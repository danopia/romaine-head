package common

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var leafUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // TODO
}

var CurrentApp *websocket.Conn

func ServeWs(handler func(Request) map[string]interface{}) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		conn, err := leafUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Client app connected")
		CurrentApp = conn

		for {
			var r Request

			if err := conn.ReadJSON(&r); err != nil {
				log.Println("Error reading app json: ", err)
				return
			}

			response := handler(r)

			conn.WriteJSON(&response)
		}
	}
}
