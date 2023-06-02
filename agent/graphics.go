package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/jpeg"
	"log"
	"os"
	"syscall"
	"unsafe"
)

var (
	libGDI32 uintptr
	token    uintptr
)

type GdiplusStartupInput struct {
	GdiplusVersion           uint32
	DebugEventCallback       uintptr
	SuppressBackgroundThread bool
	SuppressExternalCodecs   bool
}

type GdiplusStartupOutput struct {
	NotificationHook   uintptr
	NotificationUnhook uintptr
}

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

func LoadLibFunction(name string, function string) (uintptr, error) {
	lib, err := syscall.LoadLibrary(name)
	if err != nil {
		return uintptr(0), err
	}

	lib32 := uintptr(lib)
	handle, err := syscall.GetProcAddress(syscall.Handle(lib32), function)
	if err != nil {
		return uintptr(0), err
	}

	return uintptr(handle), nil
}

func GetScreen() {
	gdiPlusStartup, err := LoadLibFunction("gdiplus.dll", "GdiplusStartup")
	if err != nil {
		log.Fatal(err)
	}

	startupInput := &GdiplusStartupInput{
		GdiplusVersion: 1,
	}

	startupOutput := &GdiplusStartupOutput{}
	status, _, _ := syscall.Syscall(gdiPlusStartup, 3,
		uintptr(unsafe.Pointer(&token)),
		uintptr(unsafe.Pointer(startupInput)),
		uintptr(unsafe.Pointer(startupOutput)))

	if status != 0 {
		log.Fatalf("Gdi Syscall Error: %d", status)
	}

	getDesktopWindow, _ := LoadLibFunction("user32.dll", "GetDesktopWindow")
	hwnd, _, _ := syscall.Syscall(getDesktopWindow, 0, 0, 0, 0)
	bmp := ScreenCapture(hwnd)
	//save to memory

	size := unsafe.Sizeof(bmp)
	b := make([]byte, size)
	log.Printf("Size of bmp: %d", size)
	switch size {
	case 4:
		binary.LittleEndian.PutUint32(b, uint32(bmp))
	case 8:
		binary.LittleEndian.PutUint64(b, uint64(bmp))
	default:
		panic("uh oh")
	}

	log.Printf("Size of b: %d", len(b))

	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	var opt jpeg.Options
	opt.Quality = 100
	out, err := os.Create("./output.jpg")
	if err != nil {
		log.Fatal(err)
	}

	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		log.Fatal(err)
	}
}

func ScreenCapture(hwnd uintptr) uintptr {
	getDC, _ := LoadLibFunction("user32.dll", "GetDC")
	windowDC, _, _ := syscall.Syscall(getDC, 1, hwnd, 0, 0)
	createCompatibleDC, err := LoadLibFunction("gdi32.dll", "CreateCompatibleDC")
	if err != nil {
		log.Fatal(err)
	}

	windowCompatDC, _, _ := syscall.Syscall(createCompatibleDC, 1, hwnd, 0, 0)

	getSystemMetrics, err := LoadLibFunction("user32.dll", "GetSystemMetrics")
	if err != nil {
		log.Fatal(err)
	}

	screenx, _, _ := syscall.Syscall(getSystemMetrics, 1, uintptr(76), 0, 0)
	screeny, _, _ := syscall.Syscall(getSystemMetrics, 1, uintptr(77), 0, 0)
	width, _, _ := syscall.Syscall(getSystemMetrics, 1, uintptr(78), 0, 0)
	height, _, _ := syscall.Syscall(getSystemMetrics, 1, uintptr(79), 0, 0)

	setStretchBltMode, err := LoadLibFunction("gdi32.dll", "SetStretchBltMode")
	if err != nil {
		log.Fatal(err)
	}

	syscall.Syscall(setStretchBltMode, 2, hwnd, uintptr(1), 0)
	createCompatibleBitmap, err := LoadLibFunction("gdi32.dll", "CreateCompatibleBitmap")
	if err != nil {
		log.Fatal(err)
	}

	hbwindow, _, _ := syscall.Syscall(createCompatibleBitmap, 3, hwnd, 100, 100)

	bi := BITMAPINFOHEADER{}
	bi.BiSize = uint32(unsafe.Sizeof(bi))
	bi.BiWidth = int32(width)
	bi.BiHeight = -(int32(height))
	bi.BiPlanes = 1
	bi.BiBitCount = 32

	selectObject, err := LoadLibFunction("gdi32.dll", "SelectObject")
	if err != nil {
		log.Fatal(err)
	}

	syscall.Syscall(selectObject, 2, windowCompatDC, hbwindow, 0)

	dwBmpSize := ((int32(width)*int32(bi.BiBitCount) + 31) / 32) * 4 * int32(height)
	globalAlloc, err := LoadLibFunction("kernel32.dll", "GlobalAlloc")
	if err != nil {
		log.Fatal(err)
	}

	hDIB, _, _ := syscall.Syscall(globalAlloc, 2, 0x0042, uintptr(dwBmpSize), 0)
	lpbitmap, _, _ := syscall.Syscall(globalAlloc, 1, hDIB, 0, 0)
	stretchBlt, err := LoadLibFunction("gdi32.dll", "StretchBlt")
	if err != nil {
		log.Fatal(err)
	}

	syscall.Syscall12(stretchBlt, 11, windowCompatDC, 0, 0,
		width, height, windowDC, screenx, screeny, width, height, 0xcc0020, 0)
	getDIBits, err := LoadLibFunction("gdi32.dll", "GetDIBits")
	if err != nil {
		log.Fatal(err)
	}

	syscall.Syscall9(getDIBits, 7, windowCompatDC, hbwindow, 0, height, lpbitmap, uintptr(unsafe.Pointer(&bi)), 0, 0, 0)
	deleteDC, err := LoadLibFunction("gdi32.dll", "DeleteDC")
	if err != nil {
		log.Fatal(err)
	}

	syscall.Syscall(deleteDC, 1, windowCompatDC, 0, 0)
	releaseDC, err := LoadLibFunction("user32.dll", "ReleaseDC")
	if err != nil {
		log.Fatal(err)
	}

	syscall.Syscall(releaseDC, 1, windowDC, 0, 0)
	return hbwindow
}
