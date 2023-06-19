package main

func main() {
  go startNats()
  go connectNats()
	startWebServer()
}
