package main

import (
	"fmt"
	"log"
	"net"
	"httpfromtcp/internal/request"
)

const network = "tcp"
const port = ":42069"

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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("ERROR: Failed to read and parse: %v", err)
			continue
		}
		reqLine := req.RequestLine
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", reqLine.Method)
		fmt.Printf("- Target: %s\n", reqLine.RequestTarget)
		fmt.Printf("- Version: %s\n", reqLine.HttpVersion)
		
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}

}
