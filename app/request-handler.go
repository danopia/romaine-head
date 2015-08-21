package app

import (
	"log"

	"github.com/danopia/romaine-head/head"
)

func HandleRequest(r *Request) (response map[string]interface{}) {
	log.Printf("<<< %+v\n", r)

	response = make(map[string]interface{})
	response["context"] = r.Context

	switch r.Cmd {
	case "list chroots":
		response["chroots"] = listRoots()

	case "start chroot":
		head.StartLeaf(r.Chroot)
		response["status"] = "launching"

	case "run crouton":
		response["output"] = runCrouton(r.Args)

	case "run in chroot":
		runInChroot(r.Chroot, r.Args, r.Context)
		response["pending"] = true

	default:
		log.Printf("Client ran unknown head command " + r.Cmd)
		response["error"] = "unknown command"
	}

	log.Printf(">>> response to %s: %+v\n", r.Context, response)
	return
}
