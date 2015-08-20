package leaf

import (
	"log"

	"github.com/danopia/romaine-head/common"
)

func HandlePacket(p common.Packet) *common.Packet {
	log.Printf("head <<< %+v\n", p)

	switch p.Cmd {
	case "get info":
		//response["info"] = getVersion()

	case "execute":
		//response["output"] = runCommand(r.Args)

	default:
		log.Fatal("Head send unknown packet " + p.Cmd)
	}

	//log.Printf(">>> response to %s: %+v\n", r.Context, response)
	return nil
}
