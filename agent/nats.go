package main

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/redanthrax/gorlami/agent/message"
	"google.golang.org/protobuf/proto"
)

var (
	natsConnected bool
	natsSync      sync.Mutex
	registered    bool
	nc            *nats.Conn
	id            uuid.UUID
)

func connnectNats() {
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
		//send protobuf over nats
		msg := &message.Data{
			Category: "register",
			Message:  id.String(),
		}
		pbytes, _ := proto.Marshal(msg)
		err := nc.Publish(subject, pbytes)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Agent registered")
		nc.Subscribe(subject, func(msg *nats.Msg) {
			data := &message.Data{}
			err := proto.Unmarshal(msg.Data, data)
			if err != nil {
				log.Println(err)
			} else {
				switch data.Category {
				case "connect":
					log.Println("Do connection")
				}

				log.Printf("%#v\n", data)
			}
		})
	}

	natsSync.Unlock()
}
