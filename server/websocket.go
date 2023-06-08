package main

import (
	"log"
	"time"

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
	go func() {
		for {
			log.Println("Sending ws message")
			wsSendAgents()
			time.Sleep(5 * time.Second)
		}
	}()

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
		//message := Message{
		//	Type: "agentlist",
		//	Data: agents,
		//}

		msg := `
			<p hx-swap-oob="beforeend:#agents">asdf</p>
		`
		websocket.Message.Send(ws, msg)
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
