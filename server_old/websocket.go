package main

import (
	"encoding/json"
	"log"

	"golang.org/x/net/websocket"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func handleWebSocket(ws *websocket.Conn) {
	sessionID := getSession(ws.Request())
	if sessionID == nil {
		log.Println("Invalid session")
		ws.Close()
		return
	}

	defer ws.Close()

	log.Println("Websocket connection established")
	//do mock
	//send the agent list
	go updateWSClient(ws)
	for {
		var payload []byte
		if err := websocket.Message.Receive(ws, &payload); err != nil {
			log.Println("Error receiving message:", err)
			return
		}

		handleMessage(ws, payload)
	}
}

func handleMessage(conn *websocket.Conn, payload []byte) {
	var message Message
	if err := json.Unmarshal(payload, &message); err != nil {
		log.Println("Error parsing message:", err)
		return
	}

	log.Println(string(payload))

	switch message.Type {
	}
}

func updateWSClient(ws *websocket.Conn) {
	//do updates via websockets
	hello := []byte("hello")
	ws.Write(hello)
}
