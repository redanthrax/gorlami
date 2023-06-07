package main

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

var (
	natsConnected bool
	natsSync      sync.Mutex
	registered    bool
	nc            *nats.Conn
	id            uuid.UUID
)

func natsConnect() {
	natsSync.Lock()
	if !registered {
		id = uuid.New()
		registered = true
	}

	if !natsConnected {
		var err error
		nc, err = nats.Connect(nats.DefaultURL)
		if err != nil {
			log.Println(err)
			natsSync.Unlock()
			return
		}

		nc.SetDisconnectHandler(func(nc *nats.Conn) {
			log.Println("Disconnected...")
			natsConnected = false
			natsSync.Unlock()
		})

		subject := "agents." + id.String()
		nc.Subscribe(subject, func(msg *nats.Msg) {
			log.Printf("Received message: %s\n", string(msg.Data))

			// Send a response back to the server
			response := []byte("Message processed by agent")
			nc.Publish(msg.Reply, response)
		})

		message := []byte("Hello")
		reply, err := nc.Request("agents.server", message, 2*time.Second)
		if err != nil {
			log.Println(err)
			natsSync.Unlock()
			return
		}

		log.Printf("Received response from server: %s\n", string(reply.Data))

		natsConnected = true
	}

	natsSync.Unlock()
}
