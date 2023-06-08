package main

import "log"

type Agent struct {
	ID string
}

var (
	agents []Agent
)

func addAgent(id string) {
	exists := false
	for _, existing := range agents {
		if existing.ID == id {
			exists = true
		}
	}

	if !exists {
		newAgent := Agent{
			ID: id,
		}

		agents = append(agents, newAgent)
		log.Printf("Added agent %s, %d agents", newAgent.ID, len(agents))
	}
}
