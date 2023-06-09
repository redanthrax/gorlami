package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func startWebServer() {
	log.Println("Starting webserver")
	//serve static site
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	//websockets
	http.Handle("/ws", websocket.Handler(handleWebSocket))

	//start server
	log.Println("Listening on port 3000...")
	err := http.ListenAndServe("127.0.0.1:3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
