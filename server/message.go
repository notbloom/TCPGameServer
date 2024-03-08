package server

import (
	"fmt"
	"time"
)

// A Message contains information about the sender, the time at which the
// message was sent, and the text of the message. This gives a convenient way
// of passing the necessary information about a message from the client to the
// lobby.
type Message struct {
	time   time.Time
	client *Client
	text   string
}

// Creates a new message with the given time, client and text.
func NewMessage(time time.Time, client *Client, text string) *Message {
	return &Message{
		time:   time,
		client: client,
		text:   text,
	}
}

// Returns a string representation of the message.
func (message *Message) String() string {
	return fmt.Sprintf("%s - %s: %s\n", message.time.Format(time.Kitchen), message.client.name, message.text)
}
