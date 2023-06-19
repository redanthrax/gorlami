package main

func main() {
	go startNats()
	connectNats()
	go subscribeNats()
	startWebServer()
}
