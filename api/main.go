package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pion/webrtc/v3"
)

func main() {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)
	http.HandleFunc("/connect", connect)

	log.Print("Listening on port 3000...")
	err := http.ListenAndServe("127.0.0.1:3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func connect(w http.ResponseWriter, r *http.Request) {
	log.Println("Connect called...")

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer peerConnection.Close()

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Println(err)
		return
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		log.Println(err)
		return
	}

	offerJSON, err := json.Marshal(offer)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(offerJSON)

	log.Println("Connect end.")
}
