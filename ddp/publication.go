package ddp

import "log"
import "container/list"

type Subscription struct {
	Tube chan *Message
	Id string
}

type Publication struct {
	Name string
	Subs *list.List
	Documents map[string]map[string]interface{}
	Tube chan *Update
}

type Update struct {
	Id string
	Delete bool
	Fields map[string]interface{}
}

func CreatePublication(name string) *Publication {
	pub := &Publication{
		Name: name,
		Subs: list.New(),
		Documents: make(map[string]map[string]interface{}),
		Tube: make(chan *Update),
	}

	go func() {
		for update := range pub.Tube {
			if update.Delete {
				delete(pub.Documents, update.Id)

				pub.Broadcast(&Message{
					Type: "removed",
					Collection: pub.Name,
					Id: update.Id,
				})
			} else if update.Id != "" {

				doc, exists := pub.Documents[update.Id]

				// Update in-memory data store
				if !exists {
					doc = make(map[string]interface{})
					pub.Documents[update.Id] = doc
				}

				for key, val := range update.Fields {
					doc[key] = val
				}

				// Build and broadcast the message
				msg := &Message{
					Collection: pub.Name,
					Id: update.Id,
					Fields: update.Fields,
				}

				if exists {
					msg.Type = "changed"
				} else {
					msg.Type = "added"
				}

				pub.Broadcast(msg)
			}
		}
	}()

	return pub
}

func (pub Publication) Broadcast(msg *Message) {
	log.Println("About to broadcast")
	for e := pub.Subs.Front(); e != nil; e = e.Next() {
		e.Value.(*Subscription).Tube <- msg
		log.Println("Sent")
	}
	log.Println("Done")
}

func (pub Publication) Subscribe(sub *Subscription) {
	pub.Subs.PushBack(sub)
	log.Println("Added sub")

	for id, doc := range pub.Documents {
		sub.Tube <- &Message{
			Type: "added",
			Collection: pub.Name,
			Id: id,
			Fields: doc,
		}
	}

	sub.Tube <- &Message{
		Type: "ready",
		Subs: []string{sub.Id},
	}
}

func (pub Publication) Get(id string) map[string]interface{} {
	return pub.Documents[id]
}

func (pub Publication) Set(id string, doc map[string]interface{}) {
	pub.Tube <- &Update{
		Id: id,
		Fields: doc,
	}
}

// Helper for setting a single field
func (pub Publication) SetField(id string, key string, val interface{}) {
	fields := make(map[string]interface{})
	fields[key] = val
	pub.Set(id, fields)
}

func (pub Publication) Delete(id string) {
	pub.Tube <- &Update{
		Id: id,
		Delete: true,
	}
}
