# GoChat

A simple command-line chat application built with Go.

## Description

This application runs a TCP server that allows multiple users to connect and chat with each other. When a user connects, they are assigned a unique ID based on their remote address. Messages sent by any user are broadcast to all other connected users.

## How to Run the Server

To start the chat server, run the following command from the root of the project:

```bash
go run cmd/server/main.go
```

The server will start listening on port `4040`.

## How to Connect as a Client

You can connect to the chat server using any TCP client, such as `telnet` or `netcat`.

Open a new terminal and run:

```bash
telnet localhost 4040
```

Or:

```bash
nc localhost 4040
```

Once connected, you can start sending messages. Your messages will be seen by all other connected users, and you will see messages from them.
