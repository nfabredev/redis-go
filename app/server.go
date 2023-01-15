package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"storage"
	"strconv"
	"strings"
)

func main() {
	for {
		conn := startServer()
		handleConnection(conn)
	}
}

func startServer() net.Conn {
	fmt.Println("Starting to listen")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()
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
	go func(conn net.Conn) {
		for {
			buf := make([]byte, 1024)

			_, err := conn.Read(buf)

			if err != nil {
				if err == io.EOF {
					// connection closed by client
					break
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
	}(conn)
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

	case strings.EqualFold(command, "echo"):
		fmt.Println(len(request))
		REQUEST_LENGTH_WITH_MESSAGE := 5

		if len(request) != REQUEST_LENGTH_WITH_MESSAGE {
			return "", errors.New(encodeError("ECHO command called without message"))
		}
		response := encodeBulkString(request[3])
		return response, nil

	case strings.EqualFold(command, "set"):
		REQUEST_LENGTH_WITH_KEY := 5
		REQUEST_LENGTH_WITH_KEY_VALUE := 7

		var key string
		var value string

		if len(request) == REQUEST_LENGTH_WITH_KEY {
			key = request[3]
			value = ""
		} else if len(request) == REQUEST_LENGTH_WITH_KEY_VALUE {
			key = request[3]
			value = request[5]
		} else {
			return "", errors.New(encodeError("SET command called without key value pair to set"))
		}
		storage.Set(key, value)
		response := encodeSimpleString("OK")
		return response, nil

	case strings.EqualFold(command, "get"):
		REQUEST_LENGTH_WITH_KEY := 5

		if len(request) != REQUEST_LENGTH_WITH_KEY {
			return "", errors.New(encodeError("GET command called without a key to retrieve"))
		}
		storedKey, err := storage.Get(request[3])
		if err != nil {
			return "", errors.New(encodeError(err.Error()))
		}
		response := encodeBulkString(storedKey)
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

func encodeBulkString(string string) string {
	stringLength := strconv.Itoa(len(string))
	return "$" + stringLength + "\r\n" + string + "\r\n"
}

func handleResponse(conn net.Conn, response string) {
	conn.Write([]byte(response))
}
