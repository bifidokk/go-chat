package main

import (
	"fmt"
	"log"
)

const (
	defaultRoomName = "main"
)

type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]*Room
	receiver   chan *ClientMsg
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type ClientMsg struct {
	client *Client
	msg    []byte
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*Room),
		receiver:   make(chan *ClientMsg),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	h.rooms[defaultRoomName] = newRoom(defaultRoomName)

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case m := <-h.receiver:
			msg, err := createMessage(m.msg)

			if err != nil {
				log.Println(err)
				return
			}

			switch msg := msg.Msg.(type) {
			case *Join:
				m.client.email = msg.Email
				h.rooms[defaultRoomName].join(m.client)
			case *SendMsg:
				h.rooms[defaultRoomName].broadcast(msg)
			default:
				log.Fatalln(fmt.Sprintf("Can't resolve type of msg (%v, %T)\n", msg, msg))
			}
		}
	}
}
