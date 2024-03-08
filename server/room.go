package server

import (
	"fmt"
	"time"
)

// A ChatRoom contains the chat's name, a list of the currently connected
// clients, a history of the messages broadcast to the users in the channel,
// and the current time at which the ChatRoom will expire.
type ChatRoom struct {
	name     string
	clients  []*Client
	messages []string
	expiry   time.Time
}

// Creates an empty chat room with the given name, and sets its expiry time to
// the current time + EXPIRY_TIME.
func NewChatRoom(name string) *ChatRoom {
	return &ChatRoom{
		name:     name,
		clients:  make([]*Client, 0),
		messages: make([]string, 0),
		expiry:   time.Now().Add(EXPIRY_TIME),
	}
}

// Adds the given Client to the ChatRoom, and sends them all messages that have
// that have been sent since the creation of the ChatRoom.
func (chatRoom *ChatRoom) Join(client *Client) {
	client.chatRoom = chatRoom
	//for _, message := range chatRoom.messages {
	//	client.outgoing <- message
	//}
	chatRoom.clients = append(chatRoom.clients, client)
	chatRoom.Broadcast(fmt.Sprintf(RSP_PLAYER_JOINED, client.name))
}

// Removes the given Client from the ChatRoom.
func (chatRoom *ChatRoom) Leave(client *Client) {
	chatRoom.Broadcast(fmt.Sprintf(RSP_PLAYER_LEFT, client.name))
	for i, otherClient := range chatRoom.clients {
		if client == otherClient {
			chatRoom.clients = append(chatRoom.clients[:i], chatRoom.clients[i+1:]...)
			break
		}
	}
	client.chatRoom = nil
}

// Sends the given message to all Clients currently in the ChatRoom.
func (chatRoom *ChatRoom) Broadcast(message string) {
	chatRoom.expiry = time.Now().Add(EXPIRY_TIME)
	chatRoom.messages = append(chatRoom.messages, message)
	for _, client := range chatRoom.clients {
		client.outgoing <- message + "\n"
	}
}

// Notifies the clients within the chat room that it is being deleted, and kicks
// them back into the lobby.
func (chatRoom *ChatRoom) Delete() {
	//notify of deletion?
	chatRoom.Broadcast(NOTICE_ROOM_DELETE)
	for _, client := range chatRoom.clients {
		client.chatRoom = nil
	}
}
