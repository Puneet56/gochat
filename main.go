package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	conn    net.Conn
	message string
}

func handleMessages(channel chan Message) {
	for {
		msg := <-channel
		fmt.Fprintf(msg.conn, "Server says - %s", msg.message)
	}
}

func handleConnection(conn net.Conn, messages chan Message) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from connection: %s", err)
			conn.Close()
			break
		}

		m := string(buffer[:n])
		messages <- Message{conn, m}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatalf("Error %s", err)
	}

	channel := make(chan Message)

	go handleMessages(channel)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Error %s", err)
		}
		go handleConnection(conn, channel)
	}
}
