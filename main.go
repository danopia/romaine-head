package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", serveWs)

	log.Println("Listening on port 6205...")

	if err := http.ListenAndServe("localhost:6205", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
