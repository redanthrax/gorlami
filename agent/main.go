package main

import (
	"log"
	"sync"

	"github.com/kardianos/service"
)

const serviceName = "Agent Service"
const serviceDescription = "Agent Service for Remote Control"

var (
	serviceIsRunning bool
	programIsRunning bool
	writingSync      sync.Mutex
	libUser32        uintptr
	getCursorPos     uintptr
)

func main() {
	log.Println("Starting service...")

	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}

	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		log.Printf("Cannot create the service: %s\n", err.Error())
	}

	err = s.Run()
	if err != nil {
		log.Printf("Cannot start the service: %s\n", err.Error())
	}
}
