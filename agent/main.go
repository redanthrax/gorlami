package main

import (
	"log"
	"os"
	"sync"
	"time"
	"unsafe"

	"syscall"

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

type program struct{}

type POINT struct {
	X, Y int32
}

// Start implements service.Interface.
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

		//here be work
		lpPoint := &POINT{}
		cursorAddr, err := syscall.GetProcAddress(syscall.Handle(libUser32), "GetCursorPos")
		if err != nil {
			log.Fatal(err)
		}

		cursorPos := uintptr(cursorAddr)
		syscall.Syscall(cursorPos, 1, uintptr(unsafe.Pointer(lpPoint)), 0, 0)
		log.Printf("X: %d, Y: %d", lpPoint.X, lpPoint.Y)

		cursorAddr, err = syscall.GetProcAddress(syscall.Handle(libUser32), "SetCursorPos")
		if err != nil {
			log.Fatal(err)
		}

		//this is setting mouse position
		//cursorPos = uintptr(cursorAddr)

		//min := 10
		//max := 400
		//randGuy := rand.Intn(max-min) + min

		//syscall.Syscall(cursorPos, 2, uintptr(randGuy), uintptr(randGuy), 0)

		CaptureScreen()

		writingSync.Lock()
		programIsRunning = false
		writingSync.Unlock()
		os.Exit(0)
	}
}
