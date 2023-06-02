package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"syscall"
	"unsafe"
)

var (
	modUser32                  = syscall.NewLazyDLL("user32.dll")
	modGdi32                   = syscall.NewLazyDLL("gdi32.dll")
	procGetDC                  = modUser32.NewProc("GetDC")
	procCreateCompatibleDC     = modGdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = modGdi32.NewProc("CreateCompatibleBitmap")
	procGetDeviceCaps          = modGdi32.NewProc("GetDeviceCaps")
	procBitBlt                 = modGdi32.NewProc("BitBlt")
	procGetDIBits              = modGdi32.NewProc("GetDIBits")
	procDeleteDC               = modGdi32.NewProc("DeleteDC")
	procReleaseDC              = modUser32.NewProc("ReleaseDC")
	procSelectObject           = modGdi32.NewProc("SelectObject")
	procGetObject              = modGdi32.NewProc("GetObjectW")
)

const (
	SRCCOPY = 0x00CC0020
)

type BITMAPINFOHEADER struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type BITMAP struct {
	bmType       uint32
	bmWidth      int32
	bmHeight     int32
	bmWidthBytes int32
	bmPlanes     uint16
	bmBitsPixel  uint16
	bmBits       uintptr
}

func CaptureScreen() error {
	hDC, _, _ := procGetDC.Call(0)
	defer procReleaseDC.Call(0, hDC)

	hDest, _, _ := procCreateCompatibleDC.Call(hDC)
	defer procDeleteDC.Call(hDest)

	screenWidth := int32(GetDeviceCaps(syscall.Handle(hDC), 8))   // HORZRES
	screenHeight := int32(GetDeviceCaps(syscall.Handle(hDC), 10)) // VERTRES

	hBitmap, _, _ := procCreateCompatibleBitmap.Call(hDC, uintptr(screenWidth), uintptr(screenHeight))
	defer procDeleteDC.Call(hBitmap)

	procSelectObject.Call(hDest, hBitmap)

	_, _, _ = procBitBlt.Call(hDest, 0, 0, uintptr(screenWidth), uintptr(screenHeight), hDC, 0, 0, SRCCOPY)

	bitmapInfoHeader := BITMAPINFOHEADER{
		Size:          uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
		Width:         screenWidth,
		Height:        -screenHeight,
		Planes:        1,
		BitCount:      32,
		Compression:   0,
		SizeImage:     0,
		XPelsPerMeter: 0,
		YPelsPerMeter: 0,
		ClrUsed:       0,
		ClrImportant:  0,
	}

	var bitmap BITMAP
	_, _, _ = procGetObject.Call(hBitmap, uintptr(unsafe.Sizeof(bitmap)), uintptr(unsafe.Pointer(&bitmap)))

	bitmapSize := int(bitmap.bmWidthBytes) * int(bitmap.bmHeight)
	bits := make([]byte, bitmapSize)

	_, _, _ = procGetDIBits.Call(
		hDest,
		hBitmap,
		0,
		uintptr(screenHeight),
		uintptr(unsafe.Pointer(&bits[0])),
		uintptr(unsafe.Pointer(&bitmapInfoHeader)),
		0,
	)

	// Create an image.RGBA from the pixel data
	img := image.NewRGBA(image.Rect(0, 0, int(screenWidth), int(screenHeight)))
	stride := int(bitmap.bmWidthBytes)
	for y := 0; y < int(screenHeight); y++ {
		for x := 0; x < int(screenWidth); x++ {
			index := (y * stride) + (x * 4)
			alpha := bits[index+3]
			red := bits[index+2]
			green := bits[index+1]
			blue := bits[index]
			img.SetRGBA(x, y, color.RGBA{R: red, G: green, B: blue, A: alpha})
		}
	}

	// Save the image as PNG
	pngFile, err := os.Create("screenshot.png")
	if err != nil {
		return fmt.Errorf("Failed to create PNG file: %v", err)
	}
	defer pngFile.Close()

	err = png.Encode(pngFile, img)
	if err != nil {
		return fmt.Errorf("Failed to save screenshot as PNG: %v", err)
	}

	fmt.Println("Screenshot saved as: screenshot.png")
	return nil
}

func GetDeviceCaps(hdc syscall.Handle, index int) int {
	ret, _, _ := procGetDeviceCaps.Call(uintptr(hdc), uintptr(index))
	return int(ret)
}
