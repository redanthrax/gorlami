package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func startNatsServer() {
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
}

func connectNatsServer() {
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
}

func subscribeNatsServer() {
	nc.Subscribe("agents.*", func(msg *nats.Msg) {
		agentID := strings.TrimPrefix(msg.Subject, "agents.")
		log.Printf("Receieved message from %s: %s\n", agentID, string(msg.Data))

		response := []byte("Hello")
		nc.Publish(msg.Reply, response)
	})

	nc.SetDisconnectHandler(func(nc *nats.Conn) {
		log.Println("Disconnected...")
	})
}
