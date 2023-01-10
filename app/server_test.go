package main

import (
	"net"
	"testing"
	"time"
	"fmt"
)

func TestStartServer(t *testing.T) {
	go startServer()
	retries := 0
	for {
		_, err := net.Dial("tcp", "0.0.0.0:6379")

		if err != nil && retries > 10 {
			t.Fatal("All retries failed.")
		}

		if err != nil {
			fmt.Println("Failed to connect to port 6379, retrying in 1s")

			retries += 1
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	fmt.Println("Connection successful")

	return
}