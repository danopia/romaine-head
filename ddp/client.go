package ddp

import (
	"github.com/gorilla/websocket"
)

type ClientSub struct {
	Id          string
	Ready       bool
	Client      *Client
	Publication *Publication
}

type Client struct {
	Conn    *websocket.Conn
	Session string
	Subs    map[string]*ClientSub
	Source  chan *Message
	Sink    chan *Message
}
