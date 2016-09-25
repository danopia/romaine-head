package stem

import (
	"log"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/ddp"
	"github.com/danopia/romaine-head/head"
)

func HandlePacket(p *common.Packet, leaf *head.Leaf) {
	// log.Printf("%s <<< %+v\n", leaf.Id, p)

	switch p.Cmd {

	// Response from an execution
	case "exec":
		ddp.Commands.Set(p.Context, map[string]interface{}{
			"output": p.Extras["Output"].(string),
		})

	// TODO: wipe out all entries from this chroot
	case "collection wipe":

	case "set field":
		id := leaf.Id + "-" + p.Extras["Id"].(string)
		ddp.Apps.Set(id, map[string]interface{}{
			"Chroot":                   leaf.Id,
			p.Extras["Field"].(string): p.Extras["Value"].(map[string]interface{}),
		})
		// also has Collection

	default:
		log.Printf("Leaf sent unknown packet %s", p.Cmd)
	}

	//log.Printf(">>> response to %s: %+v\n", p.Context, response)
}
