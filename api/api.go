package api

import (
	"fmt"
	"image"
	"log"
	"net/rpc"
	"os"
	"time"

	"github.com/rickyfitts/distributed-evolution/util"
)

type EmptyMessage struct{}

type Task struct {
	Generation  int
	ID          int
	Location    image.Rectangle
	Started     time.Time
	TargetImage string
}

func GetTask() Task {
	var args EmptyMessage
	var reply Task

	Call("Master.GetTask", &args, &reply)

	return reply
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
