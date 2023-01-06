package main

import (
	"net"
	"testing"
	"fmt"
	"os/exec"
	"log"
	"time"
)

func TestMain(t *testing.T) {
	var err error
	cmd := exec.Command("./server")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	retries := 0
	for {
		_, err = net.Dial("tcp", "0.0.0.0:6379")

		if err != nil && retries > 10 {
			fmt.Println("All retries failed.")
			log.Fatal(err)
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
