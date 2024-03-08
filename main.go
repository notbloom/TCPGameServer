package main

import (
	"github.com/notbloom/wsGameServer/server"
	"log"
	"net"
	"os"
)

/*
import "github.com/notbloom/wsGameServer/server"

func main() {
	s := server.New(&server.Config{
		Host: "localhost",
		Port: "8080",
	})
	s.Run()
}
*/
// Creates a lobby, listens for client connections, and connects them to the
// lobby.
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lobby := server.NewLobby()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Println("Listening on " + server.CONN_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		lobby.Join(server.NewClient(conn))
	}
}
