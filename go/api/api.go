package api

import (
	"net/rpc"
	"os"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Job struct {
	CrossRate    float64 `json:"crossRate"`
	ID           uint32  `json:"ID"`
	MutationRate float64 `json:"mutationRate"`
	NumShapes    uint    `json:"numShapes"`
	OutputMode   string  `json:"outputMode"`
	PoolSize     uint    `json:"poolSize"`
	PopSize      uint    `json:"popSize"`
	ShapeSize    uint    `json:"shapeSize"`
	TargetImage  string  `json:"targetImage"`
}

type GetTaskArgs struct {
	WorkerID uint32
}

type Task struct {
	BestFit     eaopt.Individual `json:"bestFit"`
	Dimensions  util.Vector      `json:"dimensions"`
	Generation  uint             `json:"generation"`
	ID          int              `json:"ID"`
	Job         Job              `json:"-"`
	LastUpdate  time.Time        `json:"lastUpdate"`
	Offset      util.Vector      `json:"offset"`
	Output      string           `json:"output"`
	Status      string           `json:"status"`
	TargetImage string           `json:"-"`
	Type        string           `json:"type"`
	WorkerID    uint32           `json:"workerID"`
}

func GetTask(workerId uint32) (Task, error) {
	util.DPrintf("requesting task")

	args := GetTaskArgs{WorkerID: workerId}
	var reply Task

	err := Call("Master.GetTask", &args, &reply)

	return reply, err
}

func Update(args Task) (uint32, error) {
	util.DPrintf("sending progress")

	var reply Task

	err := Call("Master.Update", &args, &reply)

	return reply.Job.ID, err
}

// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func Call(rpcname string, args interface{}, reply interface{}) error {
	util.DPrintf("making a request")

	port := os.Getenv("PORT")
	c, err := rpc.DialHTTP("tcp", "master:"+port)
	if err != nil {
		return err
	}

	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err != nil {
		return err
	}

	return nil
}
