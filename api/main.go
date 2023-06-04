package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	server      *http.Server
	fileWatcher *fileWatcher
)

type fileWatcher struct {
	directory   string
	fileChanges chan bool
}

func (fw *fileWatcher) watch() {
	for {
		select {
		case <-fw.fileChanges:
			log.Println("File change detected. Reloading server...")
			stopServer()
			startServer()
		}
	}
}

func main() {
	fileWatcher = &fileWatcher{
		directory:   "static",
		fileChanges: make(chan bool),
	}

	err := filepath.Walk(fileWatcher.directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err = watchFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	go fileWatcher.watch()

	startServer()

	// Wait indefinitely
	select {}
}

func startServer() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(fileWatcher.directory)))

	server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()
	log.Println("Server started on http://localhost:8080")
}

func stopServer() {
	if server != nil {
		err := server.Close()
		if err != nil {
			log.Println("Error stopping server:", err)
		} else {
			log.Println("Server stopped")
		}
	}
}

func watchFile(path string) error {
	go func() {
		lastModified := getFileModTime(path)
		for {
			time.Sleep(time.Second)
			modified := getFileModTime(path)
			if modified.After(lastModified) {
				fileWatcher.fileChanges <- true
				break
			}
		}
	}()

	return nil
}

func getFileModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		log.Println("Error getting file info:", err)
		return time.Time{}
	}
	return info.ModTime()
}
