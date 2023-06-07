package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func startWebServer() {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	//websocket
	http.Handle("/ws", websocket.Handler(handleWebSocket))

	//rest api stuff
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/gorlami", gorlami)

	//start server
	go func() {
		err := http.ListenAndServe("127.0.0.1:3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("Listening on port 3000...")
	//block forever
	select {}
}

func gorlami(w http.ResponseWriter, r *http.Request) {
	log.Println("Gorlami is called")
	subject := "agents.*"
	message := []byte("this is some gorlami shit")
	nc.Publish(subject, message)
}
