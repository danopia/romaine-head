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

	Methods = make(map[string]func(args... interface{}) interface{})
}

var Methods map[string]func(args... interface{}) interface{}

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
		if handler, ok := Methods[m.Method]; ok {
			result := handler(m.Params...)

			out <- &Message{
				Type: "result",
				Id: m.Id,
				Result: result,
			}
		}
	}
}
