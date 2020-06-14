package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

func (m *Master) Echo(args *api.EchoArgs, reply *api.EchoReply) error {
	util.DPrintf("request received: { Message: %v }\n", args.Message)
	reply.Message = args.Message
	return nil
}

func (m *Master) rpcServer() {
	err := rpc.Register(m)
	if err != nil {
		log.Fatal("rpc error: ", err)
	}

	rpc.HandleHTTP()

	port := os.Getenv("RPC_PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("listener error: ", err)
	}

	log.Printf("listening for RPC on port %v\n", port)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("rpc serve error: ", err)
	}
}
