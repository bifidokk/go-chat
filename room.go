package main

type Room struct {
	name    string
	clients map[*Client]bool
}

func newRoom(name string) *Room {
	return &Room{name, make(map[*Client]bool)}
}

func (r *Room) join(c *Client) {
	r.clients[c] = true
	msg := newMessage(JoinedType, &Joined{c.email})
	c.conn.WriteJSON(msg)
}

func (r *Room) broadcast(m *SendMsg) {
	msg := newMessage(MessagesType, m)
	for c := range r.clients {
		c.conn.WriteJSON(msg)
	}
}
