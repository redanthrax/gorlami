package main

import (
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID string
}

var (
	ch chan []Agent
)

func addAgent(id string) {
	log.Printf("Adding agent %s\n", id)
	exists := false
	//for objs := range ch {
	//	for _, existing := range objs {
	//		if existing.ID == id {
	//			exists = true
	//		}
	//	}
	//}

	log.Printf("Checking if it exists")

	if !exists {
		newAgent := Agent{
			ID: id,
		}

		agentSlice := <-ch
		newAgents := append(agentSlice, newAgent)
		ch <- newAgents
		log.Printf("Added agent %s, %d agents", newAgent.ID, len(newAgents))
	} else {
		log.Printf("%s already exists", id)
	}
}

func observeAgents(ch chan []Agent) {
	previousAgents := <-ch
	for {
		log.Println("Observe")
		currentAgents := <-ch
		if !reflect.DeepEqual(previousAgents, currentAgents) {
			log.Println("Agents changed")
			//take action to notify
			previousAgents = currentAgents
		}
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
