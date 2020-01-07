package main

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	JoinType     = "join"
	SendMsgType  = "send-msg"
	JoinedType   = "joined"
	MessagesType = "msg"
)

const (
	TextMsg SendMessageType = "text"
	LogMsg                  = "log"
)

type SendMessageType string

type Message struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"msg"`
}

type Join struct {
	Email string `json:"email"`
}

type Joined struct {
	Email string `json:"email"`
}

type SendMsg struct {
	User     string          `json:"user"`
	Type     SendMessageType `json:"type"`
	Text     string          `json:"msg"`
	SendDate time.Time       `json:"time"`
}

var typeHandlers = map[string]func() interface{}{
	JoinType:    func() interface{} { return &Join{} },
	SendMsgType: func() interface{} { return &SendMsg{} },
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
