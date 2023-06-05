package main

import (
	"log"
	"syscall"
	"time"

	"github.com/kardianos/service"
)

type POINT struct {
	X, Y int32
}

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
	lib, err := syscall.LoadLibrary("user32.dll")
	if err != nil {
		log.Fatal(err)
	}

	libUser32 = uintptr(lib)
	for serviceIsRunning {
		writingSync.Lock()
		programIsRunning = true
		writingSync.Unlock()

		//getMouse()
		//img, err := CaptureScreen()
		//send img to webrtc

		writingSync.Lock()
		programIsRunning = false
		writingSync.Unlock()
	}
}
