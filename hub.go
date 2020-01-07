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
				h.joinRoom(m.client, msg)
			case *SendMsg:
				h.sendMessage(m.client, msg)
			default:
				log.Fatalln(fmt.Sprintf("Can't resolve type of msg (%v, %T)\n", msg, msg))
			}
		}
	}
}

func (h *Hub) joinRoom(client *Client, msg *Join) {
	client.room = defaultRoomName
	client.email = msg.Email
	h.rooms[defaultRoomName].join(client)
}

func (h *Hub) sendMessage(client *Client, msg *SendMsg) {
	msg.User = client.email
	h.rooms[client.room].message(msg)
}
