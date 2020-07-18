package api

import (
	"fmt"
	"net/rpc"
	"os"
	"time"
)

const RPC_TIMEOUT = 1

func Update(args TaskState) error {
	var reply Task
	return Call("Master.UpdateTask", &args, &reply)
}

// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func Call(rpcname string, args interface{}, reply interface{}) error {
	c, err := rpc.DialHTTP("tcp", os.Getenv("MASTER_URL"))
	if err != nil {
		return err
	}

	defer c.Close()

	_, timeout := handleRPC(func() bool {
		err = c.Call(rpcname, args, reply)
		return true
	})

	if timeout {
		return fmt.Errorf("rpc call timed out")
	}

	return err
}

func handleRPC(rpcCall func() bool) (bool, bool) {
	timeout := false
	success := false
	c := make(chan bool, 1)

	go func() { c <- rpcCall() }()

	select {
	case s := <-c:
		success = s
	case <-time.After(time.Second):
		timeout = true
	}

	return success, timeout
}
