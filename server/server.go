package server

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

		log.Printf("Recieved %s \n", string(buffer[:n]))

		m := string(buffer[:n])
		messages <- Message{conn, m}
	}
}

func Init() {
	ln, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatalf("Error %s", err)
	}

	log.Println("Listening on 6969")

	channel := make(chan Message)

	go handleMessages(channel)

	for {
		conn, err := ln.Accept()
		log.Printf("Connection from %s", conn.RemoteAddr())
		if err != nil {
			log.Fatalf("Error %s", err)
		}
		go handleConnection(conn, channel)
	}
}
