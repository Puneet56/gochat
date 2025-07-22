// Package server provides the server for the chat application
package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type User struct {
	ID   string
	Conn net.Conn
}

type Message struct {
	UserID string
	Msg    string
}

var (
	userChan     = make(chan User)
	brodcastChan = make(chan Message)
)

func Start() {
	wg := new(sync.WaitGroup)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go startServer(wg, ctx, ":4040")

	userList := make(map[string]User)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for u := range userChan {
			userList[u.ID] = u
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for m := range brodcastChan {
			msg := strings.TrimSpace(m.Msg)
			if msg == "" {
				continue
			}
			m.Msg = msg

			fmt.Printf("[%s]: %s\n", m.UserID, m.Msg)

			for _, u := range userList {
				if u.ID != m.UserID {
					fmt.Fprintf(u.Conn, "[%s]: %s\n", m.UserID, m.Msg)
				}
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Shutting down...")
		close(userChan)
		close(brodcastChan)
		cancel()

		<-c
		fmt.Println("Forced shutdown")
		os.Exit(0)
	}()

	wg.Wait()
	fmt.Println("Server closed")
}

func startServer(wg *sync.WaitGroup, ctx context.Context, port string) {
	defer wg.Done()
	listner, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on port " + port)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// TODO: fix shutdown bug here. Accept waits
			conn, err := listner.Accept()
			if err != nil {
				panic(err)
			}

			fmt.Println("User joined", conn.RemoteAddr().String())
			u := User{ID: conn.RemoteAddr().String(), Conn: conn}
			userChan <- u

			go handleConnection(u)
		}
	}
}

func handleConnection(user User) {
	defer user.Conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := user.Conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("User disconnected", user.Conn.RemoteAddr().String())
				return
			}
			panic(err)
		}
		brodcastChan <- Message{UserID: user.ID, Msg: string(buf[:n])}
	}
}
