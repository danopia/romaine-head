package leaf

import (
	"log"

	"github.com/danopia/romaine-head/common"
)

func HandleRequest(r common.Request) (response map[string]interface{}) {
	log.Printf("<<< %+v\n", r)

	response = make(map[string]interface{})
	response["context"] = r.Context

	switch r.Cmd {
	case "get info":
		response["info"] = getVersion()

	case "execute":
		response["port"] = runCommand(r.Args)

	default:
		log.Fatal("Client ran unknown leaf command " + r.Cmd)
	}

	log.Printf(">>> response to %s: %+v\n", r.Context, response)
	return
}
