package main

import (
	"log"
	"sync"

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

func startNats() {
	natsSync.Lock()
	if !natsConnected {
		var err error
		nc, err = nats.Connect(nats.DefaultURL)
		if err != nil {
			log.Fatal(err)
		}

		natsConnected = true

		log.Println("Connected to nats server")
	}

	natsSync.Unlock()
}

func registerNats() {
	natsSync.Lock()
	if !registered {
		id = uuid.New()
		registered = true
		subject := "agents." + id.String()
		err := nc.Publish(subject, []byte("Hello :)"))
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Agent registered")
		nc.Subscribe(subject, func(msg *nats.Msg) {
			request := string(msg.Data)
			log.Println(request)
		})
	}

	natsSync.Unlock()
}
