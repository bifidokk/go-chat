package main

const (
	JoinCommand = "join"
)

type Event struct {
	Name string `json:"event"`
}

type Message struct {
	Data MessageData `json:"data"`
}

type MessageData interface{}

type JoinMessage struct {
	Email string `json:"email"`
}

type SendMessage struct {
	Message string `json:"message"`
}

func NewMessageData(event Event) MessageData {
	switch event.Name {
	case JoinCommand:
		return JoinMessage{}
	}

	return SendMessage{}
}
