package main

import (
	"log"
	"fmt"
	"flag"
	"net/http"

	"github.com/danopia/romaine-head/head"
	"github.com/danopia/romaine-head/leaf"
	"github.com/danopia/romaine-head/common"
)

func main() {
	var port = flag.Int("port", 6205, "TCP port to listen on")
	var mode = flag.String("mode", "head", "Mode to run in (head or leaf)")
	flag.Parse()

	switch *mode {
	case "head":
		http.HandleFunc("/ws", common.ServeWs(head.HandleRequest))
		head.WaitForShutdown()
		defer head.ShutdownLeaves()

	case "leaf":
		http.HandleFunc("/ws", common.ServeWs(leaf.HandleRequest))
	}

	host := fmt.Sprint("localhost:", *port)
	log.Printf("Listening on %s...", host)
	if err := http.ListenAndServe(host, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
