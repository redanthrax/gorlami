package main

import "log"

func main() {
	log.Println("Starting Gorlami server.")
	//connect to nats
	connectNats()
	//start webserver
	startWebserver()
}
