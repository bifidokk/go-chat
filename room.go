package main

import (
	"fmt"
	"time"
)

type Room struct {
	name    string
	clients map[*Client]bool
}

func newRoom(name string) *Room {
	return &Room{name, make(map[*Client]bool)}
}

func (r *Room) join(c *Client) {
	r.clients[c] = true
	msg := newMessage(JoinedType, &Joined{c.email})
	c.conn.WriteJSON(msg)

	msg = newMessage(MessagesType, &SendMsg{
		User:     c.email,
		Type:     LogMsg,
		Text:     fmt.Sprintf("%s joined room", c.email),
		SendDate: time.Now(),
	})

	r.broadcast(msg)
}

func (r *Room) message(m *SendMsg) {
	m.SendDate = time.Now()
	msg := newMessage(MessagesType, m)
	r.broadcast(msg)
}

func (r *Room) broadcast(message *Message) {
	for c := range r.clients {
		c.conn.WriteJSON(message)
	}
}
