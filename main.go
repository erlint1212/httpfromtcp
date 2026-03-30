package main

import (
	"fmt"
	"strings"
	"os"
	"log"
	"io"
	"net"
)

const inputFilePath = "messages.txt"
const network = "tcp"
const port = ":42069"

func handleConnection(conn net.Conn) <-chan string {
	defer conn.Close()
	log.Printf("Accepted new connection from %s\n", conn.RemoteAddr())

	messages := make(chan string)
	
	go getLinesChannel(messages, conn)

	return messages
}

func getLinesChannel(messages <-chan string, conn net.Conn) {
	buffer := make([]byte, 8)
	currentLine := ""

	go func () {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("ERROR: Couldn't read bytes from file: %v", err)
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			if len(parts) == 1 {
				currentLine += parts[0]
			} else {
				for i:=0;i < len(parts) - 1;i++ {
					currentLine += parts[i]
					messages <- currentLine
				}
				currentLine = "" + parts[len(parts)-1]
			}
		}
		if currentLine != "" {
			messages <- currentLine
		}
		close(messages)
	} ()
}

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("ERROR: Couldn't open file: %v\n", err)
	}
	defer file.Close()

	listener, err := net.Listen(network, port)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect \"%s\" to port \"%s\": %v\n", network, port, err))
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("ERROR: Failed accepting connection: %v\n", err)
			continue
		}

		for msg := range handleConnection(conn) {
			fmt.Printf("read: %s\n", msg)
		}
	}


}

