package app

import (
	"log"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/ddp"
	"github.com/danopia/romaine-head/head"
)

func RefreshChroots() {
	for _, item := range listRoots() {
		ddp.Chroots.Set(item["key"].(string), map[string]interface{}{
			"status":    item["state"].(string),
			"encrypted": item["encrypted"].(bool),
			"targets":   item["targets"].([]string),
			"distro":    "precise",
		})
	}

	log.Printf("Refreshed chroot collection")
}

func init() {
	ddp.Methods["start chroot"] = func(c *ddp.Client, args ...interface{}) interface{} {
		chroot := args[0].(string)
		head.StartLeaf(chroot)
		return true
	}

	ddp.Methods["stop chroot"] = func(c *ddp.Client, args ...interface{}) interface{} {
		chroot := args[0].(string)
		if leaf, ok := head.GetLeaf(chroot); ok {
			if leaf.Conn != nil {
				leaf.State = "stopping"
				ddp.Chroots.SetField(chroot, "status", leaf.State)

				leaf.Conn.WriteJSON(&common.Packet{
					Cmd: "shutdown",
				})
				log.Printf("Issued shutdown command to %s", chroot)
				return true
			}
		}

		log.Printf("chroot %s isn't running", chroot)
		return false
	}
}

/*
	case "start chroot":
		leaf := head.StartLeaf(r.Chroot)
		leaf.PendingContext = r.Context

		response["status"] = "starting"
		response["pending"] = true

	case "run crouton":
		response["output"] = runCrouton(r.Args)

	case "build chroot":
		response["output"] = buildChroot(r.Args, r.Extras["stdin"].(string))

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
*/
