package main

import (
	"fmt"
	"strings"
	"log"
	"io"
	"net"
)

const network = "tcp"
const port = ":42069"

func getLinesChannel(conn net.Conn) <-chan string {

	messages := make(chan string)
	
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
				}
				messages <- currentLine
				currentLine = "" + parts[len(parts)-1]
			}
		}
		if currentLine != "" {
			messages <- currentLine
		}
		close(messages)
		conn.Close()
		fmt.Println("Connection closed")
	} ()

	return messages
}

func main() {
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
		fmt.Printf("Accepted new connection from %s\n", conn.RemoteAddr())

		for msg := range getLinesChannel(conn) {
			fmt.Println(msg)
		}
	}


}

