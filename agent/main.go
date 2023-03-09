package main

import (
    "fmt"
    "time"
)

func main() {
    for {
        fmt.Println("I'm the agent")
        time.Sleep(time.Second * 2)
    }
}
