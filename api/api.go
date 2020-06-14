package api

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/rickyfitts/distributed-evolution/util"
)

type EchoArgs struct {
	Message string
}

type EchoReply struct {
	Message string
}

// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func Call(rpcname string, args interface{}, reply interface{}) bool {
	util.DPrintf("making a request")

	port := os.Getenv("PORT")
	c, err := rpc.DialHTTP("tcp", "master:"+port)
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
