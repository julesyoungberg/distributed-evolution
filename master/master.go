package main

import "time"

func main() {
	for {
		println("[MASTER]: Hello World!")
		time.Sleep(10 * time.Second)
	}
}
