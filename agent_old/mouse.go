package main

import (
	"log"
	"syscall"
	"unsafe"
)

var (
	libUser32    uintptr
	getCursorPos uintptr
)

type POINT struct {
	X, Y int32
}

func getMouse() {
	lib, err := syscall.LoadLibrary("user32.dll")
	if err != nil {
		log.Fatal(err)
	}

	libUser32 = uintptr(lib)
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
}
