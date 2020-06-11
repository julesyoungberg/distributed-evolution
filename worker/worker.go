package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"
)

// TODO move EchoArgs and EchoReply to shared package
type EchoArgs struct {
	Message string
}

type EchoReply struct {
	Message string
}

type Worker struct {
	// worker state
}

// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	println("making a request")

	port := os.Getenv("MASTER_PORT")
	c, err := rpc.DialHTTP("tcp", "master:"+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func main() {
	for {
		time.Sleep(10 * time.Second)

		args := EchoArgs{Message: "Hello World!"}
		var reply EchoReply

		call("Master.Echo", &args, &reply)

		fmt.Printf("response received: { Message: %v }\n", reply.Message)
	}
}
