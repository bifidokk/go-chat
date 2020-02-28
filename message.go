package main

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	JoinType        = "join"
	SendMsgType     = "send-msg"
	GetUserListType = "get-users"
	JoinedType      = "joined"
	MessagesType    = "msg"
	UserListType    = "users"
)

const (
	TextMsg = "text"
	LogMsg  = "log"
)

type Message struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"msg"`
}

type Join struct {
	Email string `json:"email"`
}

type Joined struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type SendMsg struct {
	Email    string    `json:"email"`
	Type     string    `json:"type"`
	Text     string    `json:"msg"`
	SendDate time.Time `json:"date"`
}

type GetUserList struct{}

type UserList struct {
	Users []Joined `json:"users"`
}

var typeHandlers = map[string]func() interface{}{
	JoinType:        func() interface{} { return &Join{} },
	SendMsgType:     func() interface{} { return &SendMsg{} },
	GetUserListType: func() interface{} { return &GetUserList{} },
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
