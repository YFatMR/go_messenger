package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func echoServer(ws *websocket.Conn) {
	for {
		// Read a message from the client
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			// The connection was closed, so we'll break out of the loop
			break
		}

		// Echo the message back to the client
		err = websocket.Message.Send(ws, msg)
		if err != nil {
			// The connection was closed, so we'll break out of the loop
			break
		}
	}
}

func main() {
	http.Handle("/echo", websocket.Handler(echoServer))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
