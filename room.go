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
	msg := newMessage(JoinedType, &Joined{
		Email: c.email,
		Date:  time.Now(),
	})

	r.broadcast(msg)
}

func (r *Room) remove(c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
	}
}

func (r *Room) message(m *SendMsg) {
	m.SendDate = time.Now()
	msg := newMessage(MessagesType, m)
	r.broadcast(msg)
}

func (r *Room) userList(c *Client) {
	var users []string
	for client := range r.clients {
		users = append(users, client.email)
	}

	msg := newMessage(UserListType, &UserList{
		Users: users,
	})

	c.conn.WriteJSON(msg)
}

func (r *Room) broadcast(message *Message) {
	for c := range r.clients {
		fmt.Println(message)
		c.conn.WriteJSON(message)
	}
}
