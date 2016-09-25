package app

import (
	"log"

	"math/rand"

	"github.com/danopia/romaine-head/common"
	"github.com/danopia/romaine-head/ddp"
	"github.com/danopia/romaine-head/head"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

/*
func RefreshChroots() {
	for _, item := range listRoots() {
		ddp.Chroots.Set(item["key"].(string), map[string]interface{}{
			"status": item["state"].(string),
			"distro": "precise",
		})
	}

	log.Printf("Refreshed chroot collection")
}
*/

func init() {
	// Garbage collect a command from memory
	// Issued after client is done with the output
	ddp.Methods["/commands/expire"] = func(c *ddp.Client, args ...interface{}) interface{} {
		id := args[0].(string)
		ddp.Commands.Delete(id)
		return true
	}

	// Spool off a command to run
	ddp.Methods["/commands/exec"] = func(c *ddp.Client, args ...interface{}) interface{} {
		chroot := args[0].(string)
		path := args[1].(string)
		params := args[2].([]interface{})

		if leaf, ok := head.GetLeaf(chroot); ok {
			if leaf.Conn != nil {
				id := randSeq(20)
				// log.Printf("Running `%s %v` on %s as %s", path, params, chroot, id)

				ddp.Commands.Set(id, map[string]interface{}{
					"chroot": chroot,
					"path":   path,
					"args":   params,
					// TODO: start time
				})

				Extras := map[string]interface{}{
					"Path": path,
					"Args": params,
				}

				if len(args) > 3 {
					Extras["Stdin"] = args[3].(string)
				}

				leaf.Conn.WriteJSON(&common.Packet{
					Context: id,
					Cmd:     "exec",
					Extras:  Extras,
				})

				return id
			}
		}

		log.Printf("chroot %s isn't running", chroot)
		return nil // todo: throw error
	}
}
