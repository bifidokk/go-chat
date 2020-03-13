package main

import (
	"encoding/json"
	"strconv"
	"time"
)

// Client message type
const (
	JoinType        = "join"
	SendMsgType     = "send-msg"
	GetUserListType = "get-users"
	GetRoomListType = "get-rooms"

	JoinedType   = "joined"
	MessagesType = "msg"
	UserListType = "users"
	LeftType     = "left"
	RoomListType = "rooms"
)

// Client message content type
const (
	TextMsg = "text"
	LogMsg  = "log"
)

// Message is abstract message struct
type Message struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"msg"`
}

// Join message
type Join struct {
	Email string `json:"email"`
}

// Joined message
type Joined struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// SendMsg message
type SendMsg struct {
	Email    string    `json:"email"`
	Type     string    `json:"type"`
	Text     string    `json:"msg"`
	SendDate time.Time `json:"date"`
}

// GetUserList message
type GetUserList struct{}

// UserList message
type UserList struct {
	Users []Joined `json:"users"`
}

// Left message
type Left struct {
	Email string `json:"email"`
}

type RoomData struct {
	Name string `json:"name"`
}

// GetRoomList message
type GetRoomList struct{}

// UserList message
type RoomList struct {
	Rooms []RoomData `json:"rooms"`
}

var typeHandlers = map[string]func() interface{}{
	JoinType:        func() interface{} { return &Join{} },
	SendMsgType:     func() interface{} { return &SendMsg{} },
	GetUserListType: func() interface{} { return &GetUserList{} },
	GetRoomListType: func() interface{} { return &GetRoomList{} },
}

func createMessage(input []byte) (Message, error) {
	var raw json.RawMessage
	message := Message{
		Msg: &raw,
	}

	jsonInput, err := strconv.Unquote(string(input))
	if err != nil {
		return message, err
	}

	if err := json.Unmarshal([]byte(jsonInput), &message); err != nil {
		return message, err
	}

	msg := typeHandlers[message.Type]()
	if err := json.Unmarshal(raw, msg); err != nil {
		return message, err
	}

	message.Msg = msg

	return message, err
}

func newMessage(messageType string, message interface{}) *Message {
	return &Message{
		Type: messageType,
		Msg:  message,
	}
}
