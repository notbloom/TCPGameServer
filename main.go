package main

import "github.com/notbloom/wsGameServer/server"

func main() {
	s := server.New(&server.Config{
		Host: "localhost",
		Port: "8080",
	})
	s.Run()
}
