package ddp

import "log"

type Publication struct {
	Name string
	Subs map[string]*ClientSub
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
		Subs: make(map[string]*ClientSub),
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
	for _, sub := range pub.Subs {
		sub.Client.Sink <- msg
	}
}

func (pub Publication) Subscribe(sub *ClientSub) {
	pub.Subs[sub.Id] = sub
	sub.Client.Subs[sub.Id] = sub
	log.Println("Added sub")

	for id, doc := range pub.Documents {
		sub.Client.Sink <- &Message{
			Type: "added",
			Collection: pub.Name,
			Id: id,
			Fields: doc,
		}
	}

	sub.Client.Sink <- &Message{
		Type: "ready",
		Subs: []string{sub.Id},
	}
}

func (pub Publication) Unsubscribe(sub *ClientSub) {
	delete(pub.Subs, sub.Id)
	delete(sub.Client.Subs, sub.Id)
	log.Println("Removed sub")

	for id, _ := range pub.Documents {
		sub.Client.Sink <- &Message{
			Type: "removed",
			Collection: pub.Name,
			Id: id,
		}
	}

	sub.Client.Sink <- &Message{
		Type: "nosub",
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
