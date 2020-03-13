package main

import (
	"fmt"
	"log"
)

const (
	defaultRoomName = "main"
)

// Hub : store all chat objects
type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]*Room
	receiver   chan *ClientMsg
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// ClientMsg : message from client
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
				h.leave(client)
			}
		case m := <-h.receiver:
			msg, err := createMessage(m.msg)

			if err != nil {
				log.Println(err)
				h.leave(m.client)
				return
			}

			if !m.client.joinedRoom() && msg.Type != JoinType {
				log.Println("This user isn't joined")
				h.leave(m.client)
				return
			}

			switch msg := msg.Msg.(type) {
			case *Join:
				h.joinRoom(m.client, msg)
			case *SendMsg:
				h.sendMessage(m.client, msg)
			case *GetUserList:
				h.rooms[m.client.room].userList(m.client)
			case *GetRoomList:
				h.roomList(m.client)
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

func (h *Hub) leaveRoom(client *Client) {
	if !client.joinedRoom() {
		return
	}

	room := h.rooms[client.room]

	if room != nil {
		room.remove(client)
	}
}

func (h *Hub) sendMessage(client *Client, msg *SendMsg) {
	msg.Email = client.email
	h.rooms[client.room].message(msg)
}

func (h *Hub) leave(client *Client) {
	h.leaveRoom(client)
	delete(h.clients, client)
	close(client.send)
}

func (h *Hub) roomList(client *Client) {
	var rooms []RoomData

	for _, room := range h.rooms {
		rooms = append(rooms, RoomData{
			Name: room.name,
		})
	}

	msg := newMessage(RoomListType, &RoomList{
		Rooms: rooms,
	})

	client.conn.WriteJSON(msg)
}
