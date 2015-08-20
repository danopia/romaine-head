package head

import (
	"log"

	"github.com/danopia/romaine-head/common"
)

func HandleRequest(r common.Request) (response map[string]interface{}) {
	log.Printf("<<< %+v\n", r)

	response = make(map[string]interface{})
	response["context"] = r.Context

	switch r.Cmd {
	case "list chroots":
		response["chroots"] = listRoots()

	case "start chroot":
		//response["port"] = startLeaf(r.Chroot).Port

	case "run crouton":
		response["output"] = runCrouton(r.Args)

	case "run in chroot":
		response["output"] = runInChroot(r.Chroot, r.Args)

	default:
		log.Fatal("Client ran unknown head command " + r.Cmd)
	}

	log.Printf(">>> response to %s: %+v\n", r.Context, response)
	return
}
