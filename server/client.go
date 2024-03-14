package server

import (
	"bufio"
	"github.com/google/uuid"
	"log"
	"net"
	"strings"
	"time"
)

// A client abstracts away the idea of a connection into incoming and outgoing
// channels, and stores some information about the client's state, including
// their current name and chat room.
type Client struct {
	uid      string
	name     string
	chatRoom *Room
	incoming chan *Message
	outgoing chan string
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	isAdmin  bool
}

// Returns a new client from the given connection, and starts a reader and
// writer which receive and send information from the socket
func NewClient(conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	clientID := uuid.New()

	client := &Client{
		uid:      clientID.String(),
		name:     CLIENT_NAME,
		chatRoom: nil,
		incoming: make(chan *Message),
		outgoing: make(chan string),
		conn:     conn,
		reader:   reader,
		writer:   writer,
		isAdmin:  false,
	}

	client.Listen()
	return client
}

// Starts two threads which read from the client's outgoing channel and write to
// the client's socket connection, and read from the client's socket and write
// to the client's incoming channel.
func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

// Reads in strings from the Client's socket, formats them into Messages, and
// puts them into the Client's incoming channel.
func (client *Client) Read() {
	for {
		str, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			// TODO FIX THIS
			//if !errors.Is(err, io.EOF) {
			break
			//}
		}
		message := NewMessage(time.Now(), client, strings.TrimSuffix(str, "\n"))
		client.incoming <- message
	}
	close(client.incoming)
	log.Println("Closed client's incoming channel read thread")
}

// Reads in messages from the Client's outgoing channel, and writes them to the
// Client's socket.
func (client *Client) Write() {
	for str := range client.outgoing {
		_, err := client.writer.WriteString(str)
		if err != nil {
			log.Println(err)
			break
		}
		println("sent: " + str)
		err = client.writer.Flush()
		if err != nil {
			log.Println(err)
			break
		}
	}
	log.Println("Closed client's write thread")
}

// Closes the client's connection. Socket closing is by error checking, so this
// takes advantage of that to simplify the code and make sure all the threads
// are cleaned up.
func (client *Client) Quit() {
	client.conn.Close()
}
