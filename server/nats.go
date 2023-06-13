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

var (
	nc *nats.Conn
)

func startNats() {
	//start the actual server
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

func connectNats() {
	//hold for connection to nats server
	for {
		var err error
		nc, err = nats.Connect(nats.DefaultURL)
		if err == nil {
			break
		}

		log.Println("Waiting for NATS server...")
		time.Sleep(1 * time.Second)
	}

	log.Println("Connected to local NATS server")
}

func subscribeNats() {
	var err error
	_, err = nc.Subscribe("agents.*", func(msg *nats.Msg) {
		agentID := strings.TrimPrefix(msg.Subject, "agents.")
		addAgent(agentID)
	})
	if err != nil {
		log.Fatal(err)
	}

	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		log.Println("Disconnected")
	})
}

func monitorNatsAgents() {
	//for {
	//	log.Println("Checking nats agents")
	//	for _, agent := range <-agents {
	//		log.Println("Checking " + agent.ID)
	//		msg, err := nc.Request("agents."+agent.ID, []byte("ping"), 1*time.Second)
	//		if err != nil {
	//			log.Println(err)
	//		} else {
	//			log.Println(string(msg.Data))
	//		}
	//	}

	//	//wsSendAgent()
	//	log.Printf("Done checking nats agents, %d", len(agents))
	//	time.Sleep(5 * time.Second)
	//}
}
