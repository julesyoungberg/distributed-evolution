package api

import (
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"net/rpc"
	"os"
	"sync"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

const RPC_TIMEOUT = 1

type Output struct {
	Fitness float64
	Output  image.Image
}

type WorkerTask struct {
	BestFit     Output
	Mu          sync.Mutex
	TargetImage util.Image
	Task        Task
}

type Job struct {
	CrossRate    float64 `json:"crossRate"`
	ID           int     `json:"ID"`
	MutationRate float64 `json:"mutationRate"`
	NumShapes    uint    `json:"numShapes"`
	OverDraw     int     `json:"overDraw"`
	PoolSize     uint    `json:"poolSize"`
	PopSize      uint    `json:"popSize"`
	ShapeSize    uint    `json:"shapeSize"`
	TargetImage  string  `json:"targetImage"`
}

type Task struct {
	BestFit     eaopt.Individual  `json:"-"`
	Connected   bool              `json:"connected"`
	Dimensions  util.Vector       `json:"dimensions"`
	Generation  uint              `json:"generation"`
	ID          int               `json:"ID"`
	Job         Job               `json:"job"`
	LastUpdate  time.Time         `json:"lastUpdate"`
	Offset      util.Vector       `json:"offset"`
	Output      string            `json:"output"`
	Population  eaopt.Individuals `json:"-"`
	Status      string            `json:"status"`
	TargetImage string            `json:"targetImage"`
	Thread      int               `json:"thread"`
	Type        string            `json:"type"`
	WorkerID    uint32            `json:"workerID"`
}

func Register() {
	gob.Register(color.RGBA{})
	gob.Register(image.YCbCr{})

	gob.Register(Circle{})
	gob.Register(Polygon{})
	gob.Register(Triangle{})
	gob.Register(Shapes{})
}

func (t Task) Key() string {
	return fmt.Sprintf("task:%v", t.ID)
}

func Update(args Task) (Task, error) {
	var reply Task

	err := Call("Master.Update", &args, &reply)

	return reply, err
}

func (t Task) UpdateMaster(status string) (Task, error) {
	return Update(Task{
		ID:         t.ID,
		Generation: t.Generation,
		Job:        Job{ID: t.Job.ID},
		Status:     status,
		Thread:     t.Thread,
		WorkerID:   t.WorkerID,
	})
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
