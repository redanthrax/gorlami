package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID string
}

var (
	agents []Agent
)

func addAgent(id string) {
	log.Printf("Adding agent %s\n", id)
	exists := false
	for _, existing := range agents {
		if existing.ID == id {
			exists = true
		}
	}

	log.Printf("Checking if it exists")

	if !exists {
		newAgent := Agent{
			ID: id,
		}

		agents = append(agents, newAgent)
	}
}

func mockAddRemoveAgents() {
	for {
		//get random time in 20 seconds
		addAgent(uuid.New().String())
		randoNum := rand.Int31n(20)
		time.Sleep(time.Duration(randoNum) * time.Second)
	}
}
