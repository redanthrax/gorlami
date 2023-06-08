package main

import (
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	ws *websocket.Conn
)

func handleWebSocket(websock *websocket.Conn) {
	ws = websock
	log.Println("Websocket connection established")

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
			log.Println("sdp")
		default:
			log.Println("Unknown message type:", receivedMessage.Type)
		}
	}
}

func wsSendAgents() {
	if ws != nil {
    var sb strings.Builder
    sb.WriteString(`<div hx-swap-oob="beforeend:#agents">`)
    //load up the agent list html
    //append the div
    sb.WriteString("</div>")
		websocket.Message.Send(ws, sb.String())
	}
}

/*
func handleChatroom(chat *websocket.Conn) {
	log.Println("Handeling chats")
	for {
		var message Message
		err := websocket.JSON.Receive(chat, &message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		content := `
			<div hx-swap-oob="beforeend:#messages"><p>oi matie</p></div>`
		_ = websocket.Message.Send(chat, content)
		log.Println("Sent chataroo")
	}
}
*/
