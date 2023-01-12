// server.go

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn := startServer()
	handleConnection(conn)

	os.Exit(0)
}

func startServer() net.Conn {
	fmt.Println("Starting to listen")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Listening")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Connection accepted")

	return conn
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)

	if err != nil {
		if err == io.EOF {
			// connection closed by client
			return
		} else {
			log.Fatal("Error reading from client: ", err.Error())
		}
	}

	request, err := parseRequest(buf)
	if err != nil {
		handleResponse(conn, err.Error())
	}
	response, err := handleRequest(request)
	if err != nil {
		handleResponse(conn, err.Error())
	}
	handleResponse(conn, response)
}

func parseRequest(buf []byte) ([]string, error) {
	request := strings.Split(string(buf), "\r\n")

	fmt.Println("request:")
	fmt.Println(request)

	if !strings.HasPrefix(request[0], "*") {
		error := errors.New("The request should be a RESP Array \n" +
			strings.Join(request, " \n") +
			"See: https://redis.io/docs/reference/protocol-spec/#resp-arrays")
		return nil, error
	}
	// pop first item after above check
	_, request = request[0], request[1:]

	if !strings.HasPrefix(request[0], "$") {
		error := errors.New("Request is not bulk string \n" +
			strings.Join(request, " \n") +
			"See: https://redis.io/docs/reference/protocol-spec/#resp-bulk-strings")
		return nil, error
	}

	return request, nil
}

func handleRequest(request []string) (string, error) {
	command := request[1]

	switch {
	case strings.EqualFold(command, "ping"):
		response := encodeSimpleString("PONG")
		return response, nil
	default:
		return "", errors.New(encodeError("unknown command '" + request[1]))
	}
}

func encodeSimpleString(string string) string {
	return "+" + string + "\r\n"
}

func encodeError(string string) string {
	return "-ERR " + string + "\r\n"
}

func handleResponse(conn net.Conn, response string) {
	conn.Write([]byte(response))
}
