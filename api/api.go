package api

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/util"
)

type Task struct {
	Generation  int
	ID          int
	Location    []int
	Population  eaopt.Population
	Status      string
	TargetImage string
}

func GetTask() Task {
	var reply Task
	Call("Master.GetTask", nil, &reply)
	return reply
}

func Update(args Task) {
	Call("Master.Update", &args, nil)
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
