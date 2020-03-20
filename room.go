package main

import (
	"time"
)

// Room data
type Room struct {
	name    string
	clients map[*Client]bool
}

func newRoom(name string) *Room {
	return &Room{name, make(map[*Client]bool)}
}

func (r *Room) join(c *Client) {
	r.clients[c] = true
	c.room = r.name
	msg := newMessage(JoinedType, &Joined{
		Id:    c.id,
		Email: c.email,
		Date:  time.Now(),
		Room:  r.name,
	})

	r.broadcast(msg)
}

func (r *Room) remove(c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
	}

	msg := newMessage(LeftType, &Joined{
		Email: c.email,
	})

	r.broadcast(msg)
}

func (r *Room) message(m *SendMsg) {
	m.SendDate = time.Now()
	msg := newMessage(MessagesType, m)
	r.broadcast(msg)
}

func (r *Room) userList(c *Client) {
	var users []Joined
	for client := range r.clients {
		users = append(users, Joined{
			Email: client.email,
		})
	}

	msg := newMessage(UserListType, &UserList{
		Users: users,
	})

	c.conn.WriteJSON(msg)
}

func (r *Room) broadcast(message *Message) {
	for c := range r.clients {
		c.conn.WriteJSON(message)
	}
}
