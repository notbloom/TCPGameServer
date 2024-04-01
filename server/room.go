package server

import (
	"fmt"
	"time"
)

// A Room contains the chat's name, a list of the currently connected
// clients, a history of the messages broadcast to the users in the channel,
// and the current time at which the Room will expire.
type Room struct {
	name      string
	admin     *Client
	clients   []*Client
	messages  []string
	expiry    time.Time
	seatCount int
}

// Creates an empty chat room with the given name, and sets its expiry time to
// the current time + EXPIRY_TIME.
func NewChatRoom(name string) *Room {
	return &Room{
		name:      name,
		admin:     nil,
		clients:   make([]*Client, 0),
		messages:  make([]string, 0),
		expiry:    time.Now().Add(EXPIRY_TIME),
		seatCount: 0,
	}
}

// Adds the given Client to the Room, and sends them all messages that have
// that have been sent since the creation of the Room.
func (room *Room) Join(client *Client) {
	client.chatRoom = room
	// TODO CHECK RECONNECT and give them their old seat

	client.seat = room.seatCount
	room.seatCount++

	// SERVER HISTORY FOR RECONNECT ?
	//for _, message := range room.messages {
	//	client.outgoing <- message
	//}

	room.clients = append(room.clients, client)
	room.Broadcast(fmt.Sprintf(RSP_PLAYER_JOINED, client.name, client.seat))
}

// Removes the given Client from the Room.
func (room *Room) Leave(client *Client) {
	room.Broadcast(fmt.Sprintf(RSP_PLAYER_LEFT, client.name, client.seat))
	for i, otherClient := range room.clients {
		if client == otherClient {
			room.clients = append(room.clients[:i], room.clients[i+1:]...)
			break
		}
	}
	client.chatRoom = nil
}

// Sends the given message to all Clients currently in the Room.
func (room *Room) Broadcast(message string) {
	room.expiry = time.Now().Add(EXPIRY_TIME)
	room.messages = append(room.messages, message)
	for _, client := range room.clients {
		client.outgoing <- message + "\n"
	}
}

// Notifies the clients within the chat room that it is being deleted, and kicks
// them back into the lobby.
func (room *Room) Delete() {
	//notify of deletion?
	room.Broadcast(NOTICE_ROOM_DELETE)
	for _, client := range room.clients {
		client.chatRoom = nil
	}
}
