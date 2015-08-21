package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/danopia/romaine-head/app"
	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/head"
	"github.com/danopia/romaine-head/leaf"
	"github.com/danopia/romaine-head/stem"
)

func main() {
	var mode = flag.String("mode", "head", "Mode to run in (head or leaf)")
	var port = flag.Int("port", 6205, "TCP port that the head should be using")
	var secret = flag.String("secret", "none", "Secret token, for id and auth")
	flag.Parse()

	switch *mode {
	case "head":
		log.SetPrefix("[chronos] ")

		http.HandleFunc("/app", common.ServeWs(app.HandleConn))
		http.HandleFunc("/stem", common.ServeWs(stem.HandleLeafConn))

		head.WaitForShutdown()
		defer head.ShutdownLeaves()

		head.StartLeaf("precise")

		host := fmt.Sprint("localhost:", *port)
		log.Printf("Listening on %s...", host)
		if err := http.ListenAndServe(host, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	case "leaf":
		url := fmt.Sprint("ws://localhost:", *port, "/stem")
		log.Printf("Connecting to %s...", url)
		leaf.ConnectToHead(url, *secret)

	}
}
