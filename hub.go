package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type Hub struct {
	clients    map[*Client]bool
	rooms      map[*Room]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Room struct {
	name    string
	clients map[*Client]bool
}

func newHub() *Hub {
	room := initRoom()

	rooms := make(map[*Room]bool)
	rooms[room] = true

	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      rooms,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case data := <-h.broadcast:
			event := parseEvent(data)
			message := parseMessage(event, data)

			fmt.Println(event)
			fmt.Println(message)

			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func initRoom() *Room {
	return &Room{
		name:    "Main",
		clients: make(map[*Client]bool),
	}
}

func parseEvent(data []byte) Event {
	var event Event

	jsonInput, err := strconv.Unquote(string(data))
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal([]byte(jsonInput), &event)

	return event
}

func parseMessage(event Event, data []byte) Message {
	var message Message
	message.Data = NewMessageData(event)

	jsonInput, err := strconv.Unquote(string(data))
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal([]byte(jsonInput), &message)

	return message
}
