package main

import (
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func handleWebSocket(ws *websocket.Conn) {
	log.Println("Websocket connection established")
	//send the agent list
	sendAgents(ws)
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

func sendAgents(ws *websocket.Conn) {
	log.Println("Sending agent list")
	var sb strings.Builder
	sb.WriteString(`<div hx-swap-oob="innerHTML:#agents">`)
	for _, agent := range agents {
		sb.WriteString(`<p>` + agent.ID + `</p>`)
	}

	sb.WriteString("</div>")
	websocket.Message.Send(ws, sb.String())
}
