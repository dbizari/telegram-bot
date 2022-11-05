package main

import "fmt"

func main() {
	defer func() {
		recovery := recover()
		if recovery != nil {
			fmt.Printf("panic occured: %v", recovery)
		}
	}()

	startServer()
	startTelegramPoller()
}
