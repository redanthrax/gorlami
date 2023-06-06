package main

import (
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func main() {
	StartNats()
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/gorlami", gorlami)
	go func() {
		err := http.ListenAndServe("127.0.0.1:3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("Listening on port 3000...")
	select {}
}

func gorlami(w http.ResponseWriter, r *http.Request) {
	subject := "gorlami"
	message := []byte("this is some gorlami shit")
	nc.Publish(subject, message)
}
