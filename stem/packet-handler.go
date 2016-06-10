package stem

import (
	"log"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/head"
	"github.com/danopia/romaine-head/ddp"
)

func HandlePacket(p *common.Packet, leaf *head.Leaf) {
	// log.Printf("%s <<< %+v\n", leaf.Id, p)

	switch p.Cmd {

	// Response from an execution
	case "exec":
		ddp.Commands.Set(p.Context, map[string]interface{}{
			"output": p.Extras["Output"].(string),
		})

	default:
		log.Printf("Leaf sent unknown packet {}", p.Cmd)
	}

	//log.Printf(">>> response to %s: %+v\n", p.Context, response)
}
