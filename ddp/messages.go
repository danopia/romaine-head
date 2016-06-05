package ddp

import (
	"log"

	"github.com/danopia/romaine-head/common"
)

var Chroots *Publication
//var Apps *Publication
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
		if m.Name == "chroots" {
			Chroots.Subscribe(&Subscription{
				Tube: out,
				Id: m.Id,
			})
		} else if m.Name == "commands" {
			Commands.Subscribe(&Subscription{
				Tube: out,
				Id: m.Id,
			})
		}

	case "method":
		go runMethod(m, out)
	}
}

func runMethod(m *Message, out chan *Message) {
	if handler, ok := Methods[m.Method]; ok {
		log.Printf("Running method %s", m.Method)

		result := handler(m.Params...)
		out <- &Message{
			Type: "result",
			Id: m.Id,
			Result: result,
		}
		out <- &Message{
			Type: "updated",
			Methods: []string {m.Id},
		}

	} else {
		log.Println("Client called nonexistant DDP method", m.Method)
	}
}
