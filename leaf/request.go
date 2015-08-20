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

	case "exec":
		var oldArgs = p.Extras["Args"].([]interface{})
		args := make([]string, len(oldArgs))
		for i, v := range oldArgs {
			args[i] = v.(string)
		}

		var output = runCommand(p.Extras["Path"].(string), args)
		return &common.Packet{
			Cmd:     "exec",
			Context: p.Context,
			Extras: map[string]interface{}{
				"Output": output,
			},
		}

	default:
		log.Printf("Head sent unknown packet %s", p.Cmd)
	}

	//log.Printf(">>> response to %s: %+v\n", r.Context, response)
	return nil
}
