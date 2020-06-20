package api

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Task struct {
	BestFit     eaopt.Individual
	Generation  uint
	ID          int
	Location    []int
	Status      string
	TargetImage string
	Type        string
}

func GetTask() Task {
	util.DPrintf("requesting task")

	var args Task
	var reply Task

	Call("Master.GetTask", &args, &reply)

	return reply
}

func Update(args Task) {
	util.DPrintf("sending progress")

	var reply Task

	Call("Master.Update", &args, &reply)
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
