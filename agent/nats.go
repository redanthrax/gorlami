package main

import (
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

var (
	natsConnected bool
	natsSync      sync.Mutex
)

func natsConnect() {
	natsSync.Lock()
	if !natsConnected {
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			log.Println(err)
		}
		defer nc.Close()

		subject := "gorlami"
		_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
			log.Printf("Message: %s", string(msg.Data))
		})
		if err != nil {
			log.Println(err)
		}

		log.Printf("Subscribed to subject: %s\n", subject)
		natsConnected = true
	}
	natsSync.Unlock()
}
