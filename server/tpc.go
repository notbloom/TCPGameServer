package server

/*
import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// Server ...
type Server struct {
	host string
	port string
}

// Room ...
type Room struct {
	id      string
	name    string
	admin   *Client
	clients map[string]*Client
}

// Client ...
type Client struct {
	conn net.Conn
}

// Config ...
type Config struct {
	Host string
	Port string
}

// New ...
func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

// Run ...
func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequest()
	}
}

func (client *Client) handleRequest() {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}
		fmt.Printf("Message incoming: %s", string(message))

		// Decode message as JSON
		jsonMessage := map[string]interface{}{}
		err = json.Unmarshal([]byte(message), &jsonMessage)

		if err != nil {
			fmt.Printf("Ignoring message. Not valid JSON: %s", string(message))
			continue
		}
		if jsonMessage["type"] == nil {
			fmt.Printf("JSON has no type value: %s", string(message))
			continue
		}
		switch jsonMessage["type"] {
		case "createRoom":

			break
		// case "input" from clients, send to admin for processing
		case "input":
			jsonStructure := map[string]interface{}{
				"type":    "input",
				"content": jsonMessage["content"],
			}
			jsonOutput, err := json.Marshal(jsonStructure)
			if err != nil {

			}
			client.conn.Write(append(jsonOutput[:], []byte("\n")...))
			break
		// case "action" from admin, send to clients, validated
		case "action":
			break
		}
		//client.conn.Write([]byte("Message received." + string(message) + "\n"))
	}
}
*/
