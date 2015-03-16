package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Request represents what a client might ask
type Request struct {
	Cmd     string
	Args    []string
	Chroot  string
	Context string
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Client connected")

	for {
		var r Request

		err := conn.ReadJSON(&r)
		if err != nil {
			log.Println("Error reading json: ", err)
			return
		}
		fmt.Printf("<<< %+v\n", r)

		var response map[string]interface{}
		response = make(map[string]interface{})
		response["context"] = r.Context

		if r.Cmd == "list chroots" {
			response["chroots"] = strings.Split(listRoots(), "\n")

		} else if r.Cmd == "run crouton" {
			response["output"] = runCrouton(r.Args)

		} else if r.Cmd == "run in chroot" {
			response["output"] = runInChroot(r.Chroot, r.Args)

		} else {
			log.Fatal("Client ran unknown command " + r.Cmd)
		}

		fmt.Printf(">>> response to %s\n", r.Context)
		conn.WriteJSON(&response)
	}
}

func main() {
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe("localhost:6205", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("Listening on port 6205")
	}
}
