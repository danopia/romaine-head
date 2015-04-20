package main

import (
	"log"
	"strings"
)

// Request represents what a client might ask
type Request struct {
	Cmd     string
	Args    []string
	Chroot  string
	Context string
}

func serveRequest(r Request) (response map[string]interface{}) {
	log.Printf("<<< %+v\n", r)

	response = make(map[string]interface{})
	response["context"] = r.Context

	switch r.Cmd {
	case "list chroots":
		response["chroots"] = strings.Split(listRoots(), "\n")

	case "run crouton":
		response["output"] = runCrouton(r.Args)

	case "run in chroot":
		response["output"] = runInChroot(r.Chroot, r.Args)

	default:
		log.Fatal("Client ran unknown command " + r.Cmd)
	}

	log.Printf(">>> response to %s\n", r.Context)
	return
}
