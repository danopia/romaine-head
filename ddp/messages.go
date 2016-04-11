package ddp

import (
	"log"

	"github.com/danopia/romaine-head/common"
)

var Chroots *Publication
var Commands *Publication
func init() {
	Chroots = CreatePublication("chroots")
	Commands = CreatePublication("commands")
}

func HandleMessage(m *Message, out chan *Message) {
	log.Printf("<<< %+v", m)

	switch m.Type {

	case "connect":
		out <- &Message{
			Type: "connected",
			Session: common.GenerateSecret(),
		}

	case "ping":
		out <- &Message{
			Type: "pong",
		}

	case "sub":
		Chroots.Subscribe(&Subscription{
			Tube: out,
			Id: m.Id,
		})

	case "method":
		// TODO check Method

		Chroots.Set("sparta", map[string]interface{}{
			"status": "running",
			"distro": "trusty",
		})

		out <- &Message{
			Type: "result",
			Id: m.Id,
			Result: true,
		}
	}
}
