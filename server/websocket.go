package main

import (
	"log"

	"golang.org/x/net/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func handleWebSocket(ws *websocket.Conn) {
	defer ws.Close()

	for {
		var receivedMessage Message
		err := websocket.JSON.Receive(ws, &receivedMessage)
		if err != nil {
			log.Println("Failed to receive message:", err)
			break
		}

		log.Println("Received message:", receivedMessage)

		switch receivedMessage.Type {
		case "agents":
			log.Println("getting agents")
		case "sdp":
			log.Println("SDP:", receivedMessage.Content)
		default:
			log.Println("Unknown message type:", receivedMessage.Type)
		}
	}
}
