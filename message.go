package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	JoinType   = "join"
	JoinedType = "joined"
)

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

func createMessage(input []byte) (Message, error) {
	var msg json.RawMessage
	message := Message{
		Msg: &msg,
	}

	jsonInput, err := strconv.Unquote(string(input))
	if err != nil {
		return message, err
	}

	if err := json.Unmarshal([]byte(jsonInput), &message); err != nil {
		return message, err
	}

	switch message.Type {
	case JoinType:
		var s Join
		if err := json.Unmarshal(msg, &s); err != nil {
			return message, err
		}

		message.Msg = s
	default:
		err = errors.New(fmt.Sprintf("unknown message type: %q", message.Type))
	}

	return message, err
}

func newMessage(messageType string, message interface{}) *Message {
	return &Message{
		Type: messageType,
		Msg:  message,
	}
}
