package ddp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danopia/romaine-head/common"
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

	c := Client{
		Conn:    conn,
		Session: common.GenerateSecret(),
		Subs:    make(map[string]*ClientSub),
		Source:  make(chan *Message),
		Sink:    make(chan *Message),
	}

	log.Println("Websocket client connected")
	err = c.Conn.WriteMessage(websocket.TextMessage, []byte("o"))
	if err != nil {
		log.Println("Error writing welcome")
		return
	}

	go c.pumpMessages()
	c.Sink <- &Message{
		ServerId: "0",
	}

	go c.handleMessages()
	c.readMessages()

	for id, sub := range c.Subs {
		delete(sub.Publication.Subs, id)
	}
}

func (c *Client) pumpMessages() {
	for msg := range c.Sink {
		// log.Printf(">>> %+v", msg)

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

		writer, err := c.Conn.NextWriter(websocket.TextMessage)
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

func (c *Client) readMessages() {
	for {
		var frames []string
		if err := c.Conn.ReadJSON(&frames); err != nil {
			log.Println("Error reading websocket json: ", err)
			close(c.Sink)
			close(c.Source)
			return
		}

		for _, frame := range frames {
			var message Message
			if err := json.Unmarshal([]byte(frame), &message); err != nil {
				log.Println("Error reading websocket frame: ", err)
				return
			}
			c.Source <- &message
		}
	}
}
