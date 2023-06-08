package main

import (
	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func main() {
	//go startNats()
	//connectNats()
	//subscribeNats()
	//go checkNatsAgents()
	startWebServer()
}
