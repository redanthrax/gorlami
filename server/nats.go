package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/redanthrax/gorlami/server/message"
	"google.golang.org/protobuf/proto"
)

var (
	nc *nats.Conn
)

func connectNats() {
	//hold for connection to nats server
	for {
		var err error
		nc, err = nats.Connect("nats://dev-nats-server:4222")
		if err == nil {
			break
		}

		log.Println("Waiting for NATS server...")
		time.Sleep(1 * time.Second)
	}

	log.Println("Connected to NATS server")

	log.Println("Subscribed to agents")
	_, _ = nc.Subscribe("agents.*", func(msg *nats.Msg) {
		data := &message.Data{}
		err := proto.Unmarshal(msg.Data, data)
		if err != nil {
			log.Println(err)
		} else {
			switch data.Category {
			case "register":
				registerAgent(data.Message)
			}
		}
	})
}

func messageAgent(id string, msg *message.Data) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Println(err)
	} else {
		nc.Publish("agent."+id, data)
	}
}
