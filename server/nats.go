package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nats-io/nats.go"
)

func StartNats() {
	go func() {
		//run the nats server
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		currentDir := filepath.Dir(ex)
		server := fmt.Sprintf("%s\\nats-server.exe", currentDir)
		cmd := exec.Command(server)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	//hold for connection to nats server
	for {
		nc, err := nats.Connect(nats.DefaultURL)
		if err == nil {
			nc.Close()
			break
		}

		log.Println("Waiting for NATS server...")
		time.Sleep(1 * time.Second)
	}

	log.Println("Connected to local NATS server")

	//subscribe to gorlami
	subject := "gorlami"

	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received message: %s\n", string(msg.Data))
	})

	nc.SetDisconnectHandler(func(nc *nats.Conn) {
		log.Println("Yo we disconnected dawg")
	})

	log.Printf("Subscribed to %s", subject)
}
