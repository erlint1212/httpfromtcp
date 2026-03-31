package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"bufio"
)

const network = "udp"
const addr = "localhost:42069"

func main() {
	raddr, err := net.ResolveUDPAddr(network, addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect \"%s\" to addr \"%s\": %v\n", network, addr, err))
	}

	conn, err := net.DialUDP(network, nil, raddr)
	if err != nil {
		log.Fatalf("ERROR: Failed accepting connection: %v\n", err)
	}
	defer conn.Close()
	fmt.Printf("Accepted new connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("ERROR: Couldn't read string: %v", err)
		}

		msg_bytes := []byte(msg)
		_, err = conn.Write(msg_bytes)
		if err != nil {
			log.Printf("ERROR: Failed to writet to UDP conn: %v", err)
		}
	}


}


