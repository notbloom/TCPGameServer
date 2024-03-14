package server

import (
	"encoding/json"
	"fmt"
	uniqueid "github.com/albinj12/unique-id"
	"log"
	"strings"
	"time"
)

// A Lobby receives messages on its channels, and keeps track of the currently
// connected clients, and currently created chat rooms.
type Lobby struct {
	clients   []*Client
	chatRooms map[string]*Room
	incoming  chan *Message
	join      chan *Client
	leave     chan *Client
	delete    chan *Room
}

// Creates a lobby which beings listening over its channels.
func NewLobby() *Lobby {
	lobby := &Lobby{
		clients:   make([]*Client, 0),
		chatRooms: make(map[string]*Room),
		incoming:  make(chan *Message),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		delete:    make(chan *Room),
	}
	lobby.Listen()
	return lobby
}

// Starts a new thread which listens over the Lobby's various channels.
func (lobby *Lobby) Listen() {
	go func() {
		for {
			select {
			case message := <-lobby.incoming:
				lobby.Parse(message)
			case client := <-lobby.join:
				lobby.Join(client)
			case client := <-lobby.leave:
				lobby.Leave(client)
			case chatRoom := <-lobby.delete:
				lobby.DeleteChatRoom(chatRoom)
			}
		}
	}()
}

// Handles clients connecting to the lobby
func (lobby *Lobby) Join(client *Client) {
	if len(lobby.clients) >= MAX_CLIENTS {
		log.Println("Maximum number of clients reached")
		client.Quit()
		return
	}
	lobby.clients = append(lobby.clients, client)
	client.outgoing <- MSG_CONNECT
	log.Println("Client joined lobby")
	go func() {
		for message := range client.incoming {
			lobby.incoming <- message
		}
		lobby.leave <- client
	}()
}

// Handles clients disconnecting from the lobby.
func (lobby *Lobby) Leave(client *Client) {
	if client.chatRoom != nil {
		client.chatRoom.Leave(client)
	}
	for i, otherClient := range lobby.clients {
		if client == otherClient {
			lobby.clients = append(lobby.clients[:i], lobby.clients[i+1:]...)
			break
		}
	}
	close(client.outgoing)
	log.Println("Closed client's outgoing channel")
}

// Checks if the a channel has expired. If it has, the chat room is deleted.
// Otherwise, a signal is sent to the delete channel at its new expiry time.
func (lobby *Lobby) DeleteChatRoom(chatRoom *Room) {
	if chatRoom.expiry.After(time.Now()) {
		go func() {
			time.Sleep(chatRoom.expiry.Sub(time.Now()))
			lobby.delete <- chatRoom
		}()
		log.Println("attempted to delete chat room")
	} else {
		chatRoom.Delete()
		delete(lobby.chatRooms, chatRoom.name)
		log.Println("deleted chat room")
	}
}

// Handles messages sent to the lobby. If the message contains a command, the
// command is executed by the lobby. Otherwise, the message is sent to the
// sender's current chat room.
func (lobby *Lobby) Parse(message *Message) {

	// Decode message as JSON
	jsonMessage := map[string]interface{}{}
	err := json.Unmarshal([]byte(message.text), &jsonMessage)

	if err != nil {
		log.Println("Ignoring message. Not valid JSON: %s", message.text)

	}
	if jsonMessage["type"] == nil {
		log.Println("JSON has no type value: %s", message.text)
	}

	switch {
	default:
		lobby.SendMessage(message)
		println("Message received: " + message.text)
	case jsonMessage["type"] == TYPE_INPUT:
		lobby.SendMessage(message)
		println("Message received: " + message.text)
	case jsonMessage["type"] == TYPE_CREATE:
		name, err := uniqueid.Generateid("a", 4)
		if err != nil {
			log.Println(err)
			return
		}
		lobby.CreateChatRoom(message.client, name)
	case jsonMessage["type"] == TYPE_JOIN:
		name := jsonMessage["content"].(map[string]interface{})["code"].(string)
		fmt.Printf("%s", name)
		lobby.JoinRoom(message.client, name)
	case jsonMessage["type"] == TYPE_CONNECT:

	case strings.HasPrefix(message.text, CMD_LIST):
		lobby.ListChatRooms(message.client)
	case strings.HasPrefix(message.text, CMD_JOIN):
		name := strings.TrimSuffix(strings.TrimPrefix(message.text, CMD_JOIN+" "), "\n")
		lobby.JoinRoom(message.client, name)
	case strings.HasPrefix(message.text, CMD_LEAVE):
		lobby.LeaveChatRoom(message.client)
	case strings.HasPrefix(message.text, CMD_NAME):
		name := strings.TrimSuffix(strings.TrimPrefix(message.text, CMD_NAME+" "), "\n")
		lobby.ChangeName(message.client, name)
	case strings.HasPrefix(message.text, CMD_HELP):
		lobby.Help(message.client)
	case strings.HasPrefix(message.text, CMD_QUIT):
		message.client.Quit()
	}
}

// Attempts to send the given message to the client's current chat room. If they
// are not in a chat room, an error message is sent to the client.
func (lobby *Lobby) SendMessage(message *Message) {
	if message.client.chatRoom == nil {
		//message.client.outgoing <- ERROR_SEND
		// TODO DOBLE CHECK THIS?
		log.Println("client tried to send message in lobby")
		return
	}
	message.client.chatRoom.Broadcast(message.text)
	log.Println("client sent message" + message.text)
	//message.client.chatRoom.Broadcast(message.String())
	//log.Println("client sent message" + message.String())
}

// Attempts to create a chat room with the given name, provided that one does
// not already exist.
func (lobby *Lobby) CreateChatRoom(client *Client, name string) {
	if lobby.chatRooms[name] != nil {
		client.outgoing <- ERROR_CREATE
		log.Println("client tried to create chat room with a name already in use")
		return
	}
	chatRoom := NewChatRoom(name)
	lobby.chatRooms[name] = chatRoom
	go func() {
		time.Sleep(EXPIRY_TIME)
		lobby.delete <- chatRoom
	}()
	client.outgoing <- fmt.Sprintf(RSP_CREATE, chatRoom.name)
	log.Println("client created chat room with id: %s", name)

	// on Create Room: AutoJoin
	lobby.chatRooms[name].Join(client)
	client.isAdmin = true
	log.Println("client auto joined chat room")
}

// Attempts to add the client to the chat room with the given name, provided
// that the chat room exists.
func (lobby *Lobby) JoinRoom(client *Client, name string) {
	if lobby.chatRooms[name] == nil {
		client.outgoing <- ERROR_JOIN
		log.Println("client tried to join a chat room that does not exist")
		return
	}
	if client.chatRoom != nil {
		lobby.LeaveChatRoom(client)
	}
	lobby.chatRooms[name].Join(client)
	log.Println("client joined chat room")
}

// Removes the given client from their current chat room.
func (lobby *Lobby) LeaveChatRoom(client *Client) {
	if client.chatRoom == nil {
		client.outgoing <- ERROR_LEAVE
		log.Println("client tried to leave the lobby")
		return
	}
	client.chatRoom.Leave(client)
	log.Println("client left chat room")
}

// Changes the client's name to the given name.
func (lobby *Lobby) ChangeName(client *Client, name string) {
	if client.chatRoom == nil {
		client.outgoing <- fmt.Sprintf(NOTICE_PERSONAL_NAME, name)
	} else {
		client.chatRoom.Broadcast(fmt.Sprintf(NOTICE_ROOM_NAME, client.name, name))
	}
	client.name = name
	log.Println("client changed their name")
}

// Sends to the client the list of chat rooms currently open.
func (lobby *Lobby) ListChatRooms(client *Client) {
	client.outgoing <- "\n"
	client.outgoing <- "Chat Rooms:\n"
	for name := range lobby.chatRooms {
		client.outgoing <- fmt.Sprintf("%s\n", name)
	}
	client.outgoing <- "\n"
	log.Println("client listed chat rooms")
}

// Sends to the client the list of possible commands to the client.
func (lobby *Lobby) Help(client *Client) {
	client.outgoing <- "\n"
	client.outgoing <- "Commands:\n"
	client.outgoing <- "/help - lists all commands\n"
	client.outgoing <- "/list - lists all chat rooms\n"
	client.outgoing <- "/create foo - creates a chat room named foo\n"
	client.outgoing <- "/join foo - joins a chat room named foo\n"
	client.outgoing <- "/leave - leaves the current chat room\n"
	client.outgoing <- "/name foo - changes your name to foo\n"
	client.outgoing <- "/quit - quits the program\n"
	client.outgoing <- "\n"
	log.Println("client requested help")
}
