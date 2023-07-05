package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Agent struct {
	Name string
	ID   string
}

var (
	agents []Agent
)

func mockAgents() {
	i := 0
	for i < 10 {
		agent := Agent{
			Name: fmt.Sprintf("Agent-%d", i),
			ID:   uuid.New().String(),
		}

		agents = append(agents, agent)
		i++
	}
}

func registerAgent(id string) {
	agent := Agent{
		Name: "Agent",
		ID:   id,
	}

	agents = append(agents, agent)
}
