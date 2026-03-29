package main

import (
	"fmt"
	"strings"
	"os"
	"log"
	"io"
)

const inputFilePath = "messages.txt"

func getLinesChannel(file io.ReadCloser) <-chan string {
	buffer := make([]byte, 8)
	currentLine := ""

	messages := make(chan string)

	go func () {
		for {
			n, err := file.Read(buffer)
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
		file.Close()
	} ()

	return messages
}

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("ERROR: Couldn't open file: %v\n", err)
	}


	for msg := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", msg)
	}
}

