// server.go

package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	startServer()

	os.Exit(0)
}

func startServer() {
	fmt.Println("Starting to listen")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Listening")

	conn, err := listener.Accept()
	defer conn.Close()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Connection accepted")

	return 
}
