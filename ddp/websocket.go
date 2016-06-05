package ddp

import (
	"log"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

func ServeSockJs(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Websocket client connected")
	err = conn.WriteMessage(websocket.TextMessage, []byte("o"))
	if err != nil {
		log.Println("Error writing welcome")
		return
	}

	out := make(chan *Message, 100) // TODO: hack
	go pumpMessages(out, conn)

	out <- &Message{
		ServerId: "0",
	}

	for {
		var frames []string
		if err := conn.ReadJSON(&frames); err != nil {
			log.Println("Error reading websocket json: ", err)
			return
		}

		for _, frame := range frames {
			var message Message
			if err := json.Unmarshal([]byte(frame), &message); err != nil {
				log.Println("Error reading websocket frame: ", err)
				return
			}

			HandleMessage(&message, out)
		}
	}
}

func pumpMessages(tube chan *Message, conn *websocket.Conn) {
	for msg := range tube {
		log.Printf(">>> %+v", msg)

		frame, err := json.Marshal(msg)
		if err != nil {
			log.Println("Error serializing frame: ", err)
			return
		}

		frames := []string{string(frame)}
		payload, err := json.Marshal(frames)
		if err != nil {
			log.Println("Error serializing frames: ", err)
			return
		}

		writer, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println("Error starting to write message")
			// TODO: unsubscribe from everything
			return
		}

		writer.Write([]byte("a"))
		writer.Write(payload)
		writer.Close()
	}
}
