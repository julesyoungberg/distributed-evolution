package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

// TODO move EchoArgs and EchoReply to shared package
type EchoArgs struct {
	Message string
}

type EchoReply struct {
	Message string
}

type Master struct {
	// master state
}

func (m *Master) Echo(args *EchoArgs, reply *EchoReply) error {
	fmt.Printf("request received: { Message: %v }\n", args.Message)
	reply.Message = args.Message
	return nil
}

func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()

	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("listener error: ", err)
	}

	fmt.Printf("listening on port %v\n", port)
	http.Serve(listener, nil)
}

func main() {
	m := new(Master)
	m.server()
}
