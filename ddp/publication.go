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
	NewDoc map[string]interface{}
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
				msg := &Message{
					Collection: pub.Name,
					Id: update.Id,
					Fields: update.NewDoc,
				}

				if _, ok := pub.Documents[update.Id]; ok {
					msg.Type = "changed"
				} else {
					msg.Type = "added"
				}

				// TODO: apply changes, don't replace
				pub.Documents[update.Id] = update.NewDoc
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
		NewDoc: doc,
	}
}

func (pub Publication) Delete(id string) {
	pub.Tube <- &Update{
		Id: id,
		Delete: true,
	}
}
