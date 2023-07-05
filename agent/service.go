package main

import (
	"log"
	"sync"
	"time"

	"github.com/kardianos/service"
)

var (
	serviceIsRunning bool
	programIsRunning bool
	writingSync      sync.Mutex
)

type agent struct{}

func (a *agent) Start(s service.Service) error {
	log.Printf("%s started", s.String())
	writingSync.Lock()
	serviceIsRunning = true
	writingSync.Unlock()
	a.Run()
	return nil
}

// Stop implements service.Interface.
func (*agent) Stop(s service.Service) error {
	writingSync.Lock()
	serviceIsRunning = false
	writingSync.Unlock()
	for programIsRunning {
		log.Printf("%s stopping...", s.String())
		time.Sleep(1 * time.Second)
	}

	log.Printf("%s stopped", s.String())
	return nil
}

func (a *agent) Run() {
	for serviceIsRunning {
		writingSync.Lock()
		programIsRunning = true
		writingSync.Unlock()
		//all work here
		connnectNats()
		registerNats()
		writingSync.Lock()
		programIsRunning = false
		writingSync.Unlock()
	}
}
