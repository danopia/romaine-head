package app

import (
	"log"

	// "github.com/danopia/romaine-head/head"
	"github.com/danopia/romaine-head/ddp"
)

func RefreshChroots() {
	for _, item := range listRoots() {
		ddp.Chroots.Set(item["key"].(string), map[string]interface{}{
			"status": item["state"].(string),
			"distro": "precise",
		})
	}

	log.Printf("Refreshed chroot collection")
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
