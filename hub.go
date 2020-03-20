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
				continue
			}

			if !m.client.joinedRoom() && msg.Type != JoinType {
				log.Println("This user isn't joined")
				h.leave(m.client)
				continue
			}

			switch msg := msg.Msg.(type) {
			case *Join:
				h.join(m.client, msg)
				h.joinRoom(m.client, defaultRoomName)
			case *SendMsg:
				h.sendMessage(m.client, msg)
			case *GetUserList:
				h.rooms[m.client.room].userList(m.client)
			case *GetRoomList:
				h.roomList(m.client)
			case *AddRoom:
				h.addRoom(msg)
			case *JoinRoom:
				h.joinRoom(m.client, msg.Name)
			default:
				log.Fatalln(fmt.Sprintf("Can't resolve type of msg (%v, %T)\n", msg, msg))
			}
		}
	}
}

func (h *Hub) join(client *Client, msg *Join) {
	client.email = msg.Email
}

func (h *Hub) joinRoom(client *Client, roomName string) {
	if _, ok := h.rooms[roomName]; !ok {
		log.Printf("Room %s doesn't exist", roomName)
		return
	}

	currentRoomName := client.room

	if currentRoomName == roomName {
		log.Printf("Client %s already joined room %s", client.email, roomName)
		return
	}

	if _, ok := h.rooms[currentRoomName]; ok {
		currentRoom := h.rooms[currentRoomName]
		currentRoom.remove(client)
	}

	h.rooms[roomName].join(client)
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

func (h *Hub) addRoom(msg *AddRoom) {
	if _, ok := h.rooms[msg.Name]; ok {
		log.Printf("Room with name \"%s\" exists", msg.Name)

		return
	}

	room := newRoom(msg.Name)
	h.rooms[msg.Name] = room

	resp := newMessage(RoomAddedType, &RoomData{
		Name: room.name,
	})

	h.broadcast(resp)
}

func (h *Hub) broadcast(message *Message) {
	for c := range h.clients {
		c.conn.WriteJSON(message)
	}
}
