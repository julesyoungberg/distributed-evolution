package main

import (
	"fmt"
	"time"

	"github.com/rickyfitts/distributed-evolution/api"
)

type Worker struct {
	// worker state
}

func main() {
	for {
		time.Sleep(10 * time.Second)

		args := api.EchoArgs{Message: "Hello World!"}
		var reply api.EchoReply

		api.Call("Master.Echo", &args, &reply)

		fmt.Printf("response received: { Message: %v }\n", reply.Message)
	}
}
