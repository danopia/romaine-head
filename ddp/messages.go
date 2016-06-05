package ddp

import (
	"log"

	"github.com/danopia/romaine-head/common"
)

var Chroots *Publication
var Commands *Publication

var Methods map[string]func(args... interface{}) interface{}

func init() {
	Chroots = CreatePublication("chroots")
	Commands = CreatePublication("commands")

	Methods = make(map[string]func(args... interface{}) interface{})
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
		go runMethod(m, out)
	}
}

func runMethod(m *Message, out chan *Message) {
	if handler, ok := Methods[m.Method]; ok {
		log.Printf("Running method %s")

		result := handler(m.Params...)
		out <- &Message{
			Type: "result",
			Id: m.Id,
			Result: result,
		}

	} else {
		log.Println("Client called nonexistant DDP method", m.Method)
	}
}
