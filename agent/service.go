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

type program struct{}

func (p *program) Start(s service.Service) error {
	log.Printf("%s started", s.String())
	writingSync.Lock()
	serviceIsRunning = true
	writingSync.Unlock()
	p.Run()
	return nil
}

// Stop implements service.Interface.
func (*program) Stop(s service.Service) error {
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

func (p *program) Run() {
	for serviceIsRunning {
		writingSync.Lock()
		programIsRunning = true
		writingSync.Unlock()

		//do everything, it's a loop
		natsConnect()
		//getMouse()
		//img, err := CaptureScreen()
		//send img to webrtc

		writingSync.Lock()
		programIsRunning = false
		writingSync.Unlock()
	}
}
